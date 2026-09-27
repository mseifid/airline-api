package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"firefly-airline/internal/domains/flight"

	"gorm.io/gorm"
)

type FlightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{
		db: db,
	}
}

var _ flight.Repository = (*FlightRepository)(nil)

func (r *FlightRepository) Create(ctx context.Context, f *flight.Flight) error {
	return r.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			var airplane AirplaneModel

			if err := tx.
				First(&airplane, f.AirplaneID).
				Error; err != nil {

				if errors.Is(
					err,
					gorm.ErrRecordNotFound,
				) {
					return fmt.Errorf(
						"airplane not found",
					)
				}

				return fmt.Errorf(
					"get airplane: %w",
					err,
				)
			}

			model := FlightModel{
				AirplaneID:         f.AirplaneID,
				DepartureAirportID: f.DepartureAirportID,
				ArrivalAirportID:   f.ArrivalAirportID,
				DepartureAt:        f.DepartureAt,
				ArrivalAt:          f.ArrivalAt,
				Price:              f.Price,
				Capacity:           int(airplane.Capacity),
				AvailableSeats:     int(airplane.Capacity),
				Status:             f.Status,
			}

			if err := tx.Create(&model).Error; err != nil {
				return fmt.Errorf(
					"create flight: %w",
					err,
				)
			}

			f.ID = model.ID
			f.Capacity = model.Capacity
			f.AvailableSeats = model.AvailableSeats
			f.CreatedAt = model.CreatedAt
			f.UpdatedAt = model.UpdatedAt

			return nil
		},
	)
}

func (r *FlightRepository) GetByID(ctx context.Context, id int64) (*flight.Flight, error) {
	var model FlightModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, flight.ErrNotFound
		}

		return nil, fmt.Errorf(
			"get flight: %w",
			err,
		)
	}

	return toDomainFlight(&model), nil
}

func (r *FlightRepository) Update(ctx context.Context, f *flight.Flight) error {
	var model FlightModel

	err := r.db.WithContext(ctx).
		First(&model, f.ID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return flight.ErrNotFound
		}

		return fmt.Errorf(
			"find flight: %w",
			err,
		)
	}

	model.AirplaneID = f.AirplaneID
	model.DepartureAirportID = f.DepartureAirportID
	model.ArrivalAirportID = f.ArrivalAirportID
	model.DepartureAt = f.DepartureAt
	model.ArrivalAt = f.ArrivalAt
	model.Price = f.Price
	model.Status = f.Status

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return fmt.Errorf(
			"update flight: %w",
			err,
		)
	}

	f.Capacity = model.Capacity
	f.AvailableSeats = model.AvailableSeats
	f.CreatedAt = model.CreatedAt
	f.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *FlightRepository) ListActiveByCitiesAndDateRange(
	ctx context.Context,
	sourceCityID int64,
	destinationCityID int64,
	from time.Time,
	to time.Time,
) ([]flight.Flight, error) {
	var models []FlightModel

	err := r.db.WithContext(ctx).
		Table("flight AS f").
		Joins(
			"JOIN airport AS departure_airport ON "+
				"departure_airport.id = f.departure_airport_id",
		).
		Joins(
			"JOIN airport AS arrival_airport ON "+
				"arrival_airport.id = f.arrival_airport_id",
		).
		Where(
			"departure_airport.city_id = ?",
			sourceCityID,
		).
		Where(
			"arrival_airport.city_id = ?",
			destinationCityID,
		).
		Where(
			"f.status = ?",
			flight.StatusActive,
		).
		Where(
			"f.departure_at >= ?",
			from,
		).
		Where(
			"f.departure_at < ?",
			to,
		).
		Order("f.departure_at ASC").
		Find(&models).
		Error

	if err != nil {
		return nil, fmt.Errorf(
			"search flights: %w",
			err,
		)
	}

	result := make(
		[]flight.Flight,
		0,
		len(models),
	)

	for i := range models {
		result = append(
			result,
			*toDomainFlight(&models[i]),
		)
	}

	return result, nil
}

func (r *FlightRepository) ListActive(ctx context.Context) ([]flight.Flight, error) {
	var models []FlightModel

	err := r.db.WithContext(ctx).
		Where("status = ?", flight.StatusActive).
		Order("departure_at ASC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	flights := make([]flight.Flight, 0, len(models))

	for i := range models {
		flights = append(flights, *toDomainFlight(&models[i]))
	}

	return flights, nil
}

func toDomainFlight(model *FlightModel) *flight.Flight {
	return &flight.Flight{
		ID:                 model.ID,
		AirplaneID:         model.AirplaneID,
		DepartureAirportID: model.DepartureAirportID,
		ArrivalAirportID:   model.ArrivalAirportID,
		DepartureAt:        model.DepartureAt,
		ArrivalAt:          model.ArrivalAt,
		Price:              model.Price,
		Capacity:           model.Capacity,
		AvailableSeats:     model.AvailableSeats,
		Status:             model.Status,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}
