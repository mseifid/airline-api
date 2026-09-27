package flight

import (
	"context"
	"fmt"
	"time"
)

const (
	StatusActive    = "active"
	StatusCancelled = "cancelled"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	airplaneID int64,
	departureAirportID int64,
	arrivalAirportID int64,
	departureAt time.Time,
	arrivalAt time.Time,
	price int64,
) (Flight, error) {
	if airplaneID <= 0 {
		return Flight{}, fmt.Errorf("invalid airplane id")
	}

	if departureAirportID <= 0 {
		return Flight{}, fmt.Errorf("invalid departure airport id")
	}

	if arrivalAirportID <= 0 {
		return Flight{}, fmt.Errorf("invalid arrival airport id")
	}

	if departureAirportID == arrivalAirportID {
		return Flight{}, ErrInvalidAirports
	}

	if !arrivalAt.After(departureAt) {
		return Flight{}, ErrInvalidSchedule
	}

	if price < 0 {
		return Flight{}, ErrInvalidPrice
	}

	// Capacity is populated by the repository from the airplane.
	flight := &Flight{
		AirplaneID:         airplaneID,
		DepartureAirportID: departureAirportID,
		ArrivalAirportID:   arrivalAirportID,
		DepartureAt:        departureAt,
		ArrivalAt:          arrivalAt,
		Price:              price,
		Status:             StatusActive,
	}

	if err := s.repository.Create(ctx, flight); err != nil {
		return Flight{}, fmt.Errorf("create flight: %w", err)
	}

	return *flight, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	airplaneID int64,
	departureAirportID int64,
	arrivalAirportID int64,
	departureAt time.Time,
	arrivalAt time.Time,
	price int64,
	status string,
) (Flight, error) {
	if id <= 0 {
		return Flight{}, fmt.Errorf("invalid flight id")
	}

	if airplaneID <= 0 {
		return Flight{}, fmt.Errorf("invalid airplane id")
	}

	if departureAirportID <= 0 {
		return Flight{}, fmt.Errorf("invalid departure airport id")
	}

	if arrivalAirportID <= 0 {
		return Flight{}, fmt.Errorf("invalid arrival airport id")
	}

	if departureAirportID == arrivalAirportID {
		return Flight{}, ErrInvalidAirports
	}

	if !arrivalAt.After(departureAt) {
		return Flight{}, ErrInvalidSchedule
	}

	if price < 0 {
		return Flight{}, ErrInvalidPrice
	}

	if status == "" {
		return Flight{}, fmt.Errorf("status is required")
	}

	flight := &Flight{
		ID:                 id,
		AirplaneID:         airplaneID,
		DepartureAirportID: departureAirportID,
		ArrivalAirportID:   arrivalAirportID,
		DepartureAt:        departureAt,
		ArrivalAt:          arrivalAt,
		Price:              price,
		Status:             status,
	}

	if err := s.repository.Update(ctx, flight); err != nil {
		return Flight{}, fmt.Errorf("update flight: %w", err)
	}

	return *flight, nil
}

func (s *Service) Search(
	ctx context.Context,
	sourceCityID int64,
	destinationCityID int64,
	from time.Time,
	to time.Time,
) ([]Flight, error) {
	if sourceCityID <= 0 {
		return nil, fmt.Errorf("invalid source city id")
	}

	if destinationCityID <= 0 {
		return nil, fmt.Errorf("invalid destination city id")
	}

	if sourceCityID == destinationCityID {
		return nil, fmt.Errorf(
			"source and destination cities must be different",
		)
	}

	if to.Before(from) {
		return nil, fmt.Errorf(
			"arrival date must be on or after departure date",
		)
	}

	flights, err := s.repository.ListActiveByCitiesAndDateRange(
		ctx,
		sourceCityID,
		destinationCityID,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("search flights: %w", err)
	}

	return flights, nil
}

func (s *Service) ListActive(ctx context.Context) ([]Flight, error) {
	return s.repository.ListActive(ctx)
}