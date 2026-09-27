package ticket

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Purchase(ctx context.Context, agencyID int64, flightID int64, passengers []Passenger) (*Ticket, error) {
	if agencyID <= 0 {
		return nil, errors.New("invalid agency")
	}

	if flightID <= 0 {
		return nil, errors.New("invalid flight")
	}

	if len(passengers) == 0 {
		return nil, ErrInvalidPassengerCount
	}

	return s.repository.Purchase(
		ctx,
		agencyID,
		flightID,
		passengers,
	)
}

func (s *Service) ListByFlight(ctx context.Context, agencyID int64, flightID int64) ([]Ticket, error) {
	if agencyID <= 0 {
		return nil, errors.New("invalid agency")
	}

	if flightID <= 0 {
		return nil, errors.New("invalid flight")
	}

	return s.repository.ListByFlightAndAgency(
		ctx,
		flightID,
		agencyID,
	)
}

func (s *Service) Cancel(
	ctx context.Context,
	agencyID int64,
	ticketID int64,
) (*Ticket, error) {
	if agencyID <= 0 {
		return nil, errors.New("invalid agency")
	}

	if ticketID <= 0 {
		return nil, errors.New("invalid ticket")
	}

	t, err := s.repository.GetForCancellation(
		ctx,
		ticketID,
		agencyID,
	)
	if err != nil {
		return nil, err
	}

	if t.Status == StatusCancelled {
		return nil, ErrTicketAlreadyCancelled
	}

	now := time.Now()

	remaining := t.DepartureAt.Sub(now)

	if remaining <= 0 {
		return nil, ErrFlightAlreadyDeparted
	}

	var refundAmount int64

	switch {
	case remaining >= 4*time.Hour:
		refundAmount = t.TotalPrice * 70 / 100

	case remaining >= 1*time.Hour:
		refundAmount = t.TotalPrice * 50 / 100

	default:
		refundAmount = t.TotalPrice * 10 / 100
	}

	return s.repository.Cancel(
		ctx,
		ticketID,
		agencyID,
		refundAmount,
		now,
	)
}

func (s *Service) GetByID(
	ctx context.Context,
	ticketID int64,
	agencyID int64,
) (*Ticket, error) {
	if agencyID <= 0 {
		return nil, errors.New("invalid agency")
	}

	if ticketID <= 0 {
		return nil, errors.New("invalid ticket")
	}

	return s.repository.GetByID(
		ctx,
		ticketID,
		agencyID,
	)
}