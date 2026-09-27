package database

import (
	"context"
	"errors"
	"time"

	"firefly-airline/internal/domains/ticket"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) Purchase(
	ctx context.Context,
	agencyID int64,
	flightID int64,
	passengers []ticket.Passenger,
) (*ticket.Ticket, error) {
	var result *ticket.Ticket

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock agency first
		var agencyModel AgencyModel

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&agencyModel, agencyID).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("agency not found")
		}

		if err != nil {
			return err
		}

		// Lock flight second
		var flightModel FlightModel

		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&flightModel, flightID).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ticket.ErrFlightNotFound
		}

		if err != nil {
			return err
		}

		if flightModel.Status != ticket.StatusActive {
			return errors.New("flight is not active")
		}

		var airplaneModel AirplaneModel

		err = tx.
			Select("id", "type").
			First(&airplaneModel, flightModel.AirplaneID).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("airplane not found")
		}

		if err != nil {
			return err
		}

		seatCount := len(passengers)

		if flightModel.AvailableSeats < seatCount {
			return ticket.ErrNotEnoughSeats
		}

		// Calculate the price from the locked flight
		totalPrice := flightModel.Price * int64(seatCount)

		// Check the locked wallet
		if agencyModel.WalletBalance < totalPrice {
			return ticket.ErrInsufficientBalance
		}

		// Decrease available seats
		flightModel.AvailableSeats -= seatCount

		if err := tx.Save(&flightModel).Error; err != nil {
			return err
		}

		// Deduact wallet balance
		agencyModel.WalletBalance -= totalPrice

		if err := tx.Save(&agencyModel).Error; err != nil {
			return err
		}

		ticketModel := TicketModel{
			AgencyID:   agencyID,
			FlightID:   flightID,
			SeatCount:  seatCount,
			UnitPrice:  flightModel.Price,
			TotalPrice: totalPrice,
			Status:     ticket.StatusActive,
		}

		if err := tx.Create(&ticketModel).Error; err != nil {
			return err
		}

		// Create passengers.
		passengerModels := make([]PassengerModel, 0, len(passengers))

		for _, passenger := range passengers {
			passengerModels = append(passengerModels, PassengerModel{
				TicketID:     ticketModel.ID,
				Name:         passenger.Name,
				Mobile:       passenger.Mobile,
				NationalCode: passenger.NationalCode,
			})
		}

		if err := tx.Create(&passengerModels).Error; err != nil {
			return err
		}

		// Build domain response
		result = &ticket.Ticket{
			ID:           ticketModel.ID,
			AgencyID:     ticketModel.AgencyID,
			FlightID:     ticketModel.FlightID,
			SeatCount:    ticketModel.SeatCount,
			UnitPrice:    ticketModel.UnitPrice,
			TotalPrice:   ticketModel.TotalPrice,
			Status:       ticketModel.Status,
			CreatedAt:    ticketModel.CreatedAt,
			DepartureAt:  flightModel.DepartureAt,
			ArrivalAt:    flightModel.ArrivalAt,
			AirplaneType: airplaneModel.Type,
			Passengers:   make([]ticket.Passenger, 0, len(passengerModels)),
		}

		for _, model := range passengerModels {
			result.Passengers = append(
				result.Passengers,
				ticket.Passenger{
					ID:           model.ID,
					TicketID:     model.TicketID,
					Name:         model.Name,
					Mobile:       model.Mobile,
					NationalCode: model.NationalCode,
				},
			)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TicketRepository) ListByFlightAndAgency(ctx context.Context, flightID int64, agencyID int64) ([]ticket.Ticket, error) {
	var ticketModels []TicketModel

	err := r.db.WithContext(ctx).
		Where("flight_id = ?", flightID).
		Where("agency_id = ?", agencyID).
		Order("created_at DESC").
		Find(&ticketModels).
		Error

	if err != nil {
		return nil, err
	}

	if len(ticketModels) == 0 {
		return []ticket.Ticket{}, nil
	}

	ticketIDs := make([]int64, 0, len(ticketModels))

	for _, model := range ticketModels {
		ticketIDs = append(ticketIDs, model.ID)
	}

	var passengerModels []PassengerModel

	err = r.db.WithContext(ctx).
		Where("ticket_id IN ?", ticketIDs).
		Order("id ASC").
		Find(&passengerModels).
		Error

	if err != nil {
		return nil, err
	}

	passengersByTicket := make(map[int64][]ticket.Passenger)

	for _, model := range passengerModels {
		passengersByTicket[model.TicketID] = append(
			passengersByTicket[model.TicketID],
			ticket.Passenger{
				ID:           model.ID,
				TicketID:     model.TicketID,
				Name:         model.Name,
				Mobile:       model.Mobile,
				NationalCode: model.NationalCode,
			},
		)
	}

	result := make([]ticket.Ticket, 0, len(ticketModels))

	for _, model := range ticketModels {
		result = append(result, ticket.Ticket{
			ID:          model.ID,
			AgencyID:    model.AgencyID,
			FlightID:    model.FlightID,
			SeatCount:   model.SeatCount,
			UnitPrice:   model.UnitPrice,
			TotalPrice:  model.TotalPrice,
			Status:      model.Status,
			CreatedAt:   model.CreatedAt,
			CancelledAt: model.CancelledAt,
			Passengers:  passengersByTicket[model.ID],
		})
	}

	return result, nil
}

func (r *TicketRepository) Cancel(
	ctx context.Context,
	ticketID int64,
	agencyID int64,
	refundAmount int64,
	now time.Time,
) (*ticket.Ticket, error) {
	var result *ticket.Ticket

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Agency
		var agencyModel AgencyModel

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&agencyModel, agencyID).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("agency not found")
			}

			return err
		}

		// Find flight ID without locking ticket.
		var ticketInfo TicketModel

		if err := tx.
			Select("id", "flight_id").
			Where("id = ?", ticketID).
			Where("agency_id = ?", agencyID).
			First(&ticketInfo).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ticket.ErrTicketNotFound
			}

			return err
		}

		// Flight
		var flightModel FlightModel

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&flightModel, ticketInfo.FlightID).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ticket.ErrFlightNotFound
			}

			return err
		}

		// Ticket
		var ticketModel TicketModel

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", ticketID).
			Where("agency_id = ?", agencyID).
			First(&ticketModel).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ticket.ErrTicketNotFound
			}

			return err
		}

		if ticketModel.Status == ticket.StatusCancelled {
			return ticket.ErrTicketAlreadyCancelled
		}

		flightModel.AvailableSeats += ticketModel.SeatCount

		if flightModel.AvailableSeats > flightModel.Capacity {
			return errors.New("flight available seats exceed capacity")
		}

		if err := tx.Save(&flightModel).Error; err != nil {
			return err
		}

		agencyModel.WalletBalance += refundAmount

		if err := tx.Save(&agencyModel).Error; err != nil {
			return err
		}

		cancelledAt := now

		ticketModel.Status = ticket.StatusCancelled
		ticketModel.CancelledAt = &cancelledAt

		if err := tx.Save(&ticketModel).Error; err != nil {
			return err
		}

		result = &ticket.Ticket{
			ID:          ticketModel.ID,
			AgencyID:    ticketModel.AgencyID,
			FlightID:    ticketModel.FlightID,
			SeatCount:   ticketModel.SeatCount,
			UnitPrice:   ticketModel.UnitPrice,
			TotalPrice:  ticketModel.TotalPrice,
			Status:      ticketModel.Status,
			CreatedAt:   ticketModel.CreatedAt,
			CancelledAt: ticketModel.CancelledAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TicketRepository) GetForCancellation(ctx context.Context, ticketID int64, agencyID int64) (*ticket.Ticket, error) {
	var ticketModel TicketModel

	err := r.db.WithContext(ctx).
		Where("id = ?", ticketID).
		Where("agency_id = ?", agencyID).
		First(&ticketModel).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrTicketNotFound
	}

	if err != nil {
		return nil, err
	}

	var flightModel FlightModel

	err = r.db.WithContext(ctx).
		Select("departure_at").
		Where("id = ?", ticketModel.FlightID).
		First(&flightModel).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrFlightNotFound
	}

	if err != nil {
		return nil, err
	}

	return &ticket.Ticket{
		ID:          ticketModel.ID,
		AgencyID:    ticketModel.AgencyID,
		FlightID:    ticketModel.FlightID,
		SeatCount:   ticketModel.SeatCount,
		UnitPrice:   ticketModel.UnitPrice,
		TotalPrice:  ticketModel.TotalPrice,
		Status:      ticketModel.Status,
		CreatedAt:   ticketModel.CreatedAt,
		CancelledAt: ticketModel.CancelledAt,
		DepartureAt: flightModel.DepartureAt,
	}, nil
}

func (r *TicketRepository) GetByID(ctx context.Context, ticketID int64, agencyID int64) (*ticket.Ticket, error) {
	var ticketModel TicketModel

	err := r.db.WithContext(ctx).
		Where("id = ?", ticketID).
		Where("agency_id = ?", agencyID).
		First(&ticketModel).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrTicketNotFound
	}

	if err != nil {
		return nil, err
	}

	var passengerModels []PassengerModel

	err = r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketModel.ID).
		Order("id ASC").
		Find(&passengerModels).
		Error

	if err != nil {
		return nil, err
	}

	passengers := make([]ticket.Passenger, 0, len(passengerModels))

	for _, model := range passengerModels {
		passengers = append(passengers, ticket.Passenger{
			ID:           model.ID,
			TicketID:     model.TicketID,
			Name:         model.Name,
			Mobile:       model.Mobile,
			NationalCode: model.NationalCode,
		})
	}

	var flightModel FlightModel

	err = r.db.WithContext(ctx).
		Where("id = ?", ticketModel.FlightID).
		First(&flightModel).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ticket.ErrFlightNotFound
	}

	if err != nil {
		return nil, err
	}

	var airplaneModel AirplaneModel

	err = r.db.WithContext(ctx).
		Select("id", "type").
		Where("id = ?", flightModel.AirplaneID).
		First(&airplaneModel).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("airplane not found")
	}

	if err != nil {
		return nil, err
	}

	return &ticket.Ticket{
		ID:           ticketModel.ID,
		AgencyID:     ticketModel.AgencyID,
		FlightID:     ticketModel.FlightID,
		SeatCount:    ticketModel.SeatCount,
		UnitPrice:    ticketModel.UnitPrice,
		TotalPrice:   ticketModel.TotalPrice,
		Status:       ticketModel.Status,
		CreatedAt:    ticketModel.CreatedAt,
		CancelledAt:  ticketModel.CancelledAt,
		DepartureAt:  flightModel.DepartureAt,
		ArrivalAt:    flightModel.ArrivalAt,
		AirplaneType: airplaneModel.Type,
		Passengers:   passengers,
	}, nil
}
