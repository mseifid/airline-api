package airport

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"firefly-airline/internal/domains/city"
)

var ErrNotFound = errors.New("airport not found")
var ErrCodeAlreadyExists = errors.New("airport code already exists")

type Service struct {
	repository     Repository
	cityRepository city.Repository
}

func NewService(
	repository Repository,
	cityRepository city.Repository,
) *Service {
	return &Service{
		repository:     repository,
		cityRepository: cityRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	cityID int64,
	code string,
	name string,
) (Airport, error) {
	if cityID <= 0 {
		return Airport{}, fmt.Errorf("invalid city id")
	}

	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)

	if code == "" {
		return Airport{}, fmt.Errorf("airport code is required")
	}

	if name == "" {
		return Airport{}, fmt.Errorf("airport name is required")
	}

	if _, err := s.cityRepository.GetByID(ctx, cityID); err != nil {
		return Airport{}, fmt.Errorf("city: %w", err)
	}

	airport := &Airport{
		CityID: cityID,
		Code:   code,
		Name:   name,
	}

	if err := s.repository.Create(ctx, airport); err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			return Airport{}, ErrCodeAlreadyExists
		}

		return Airport{}, fmt.Errorf("create airport: %w", err)
	}

	return *airport, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (Airport, error) {
	airport, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Airport{}, ErrNotFound
		}

		return Airport{}, fmt.Errorf("get airport: %w", err)
	}

	return *airport, nil
}

func (s *Service) List(
	ctx context.Context,
	cityID *int64,
) ([]Airport, error) {
	if cityID != nil && *cityID <= 0 {
		return nil, fmt.Errorf("invalid city id")
	}

	airports, err := s.repository.List(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("list airports: %w", err)
	}

	return airports, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	cityID int64,
	code string,
	name string,
) (Airport, error) {
	if cityID <= 0 {
		return Airport{}, fmt.Errorf("invalid city id")
	}

	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)

	if code == "" {
		return Airport{}, fmt.Errorf("airport code is required")
	}

	if name == "" {
		return Airport{}, fmt.Errorf("airport name is required")
	}

	if _, err := s.cityRepository.GetByID(ctx, cityID); err != nil {
		return Airport{}, fmt.Errorf("city: %w", err)
	}

	airport := &Airport{
		ID:     id,
		CityID: cityID,
		Code:   code,
		Name:   name,
	}

	if err := s.repository.Update(ctx, airport); err != nil {
		if errors.Is(err, ErrNotFound) {
			return Airport{}, ErrNotFound
		}

		if errors.Is(err, ErrCodeAlreadyExists) {
			return Airport{}, ErrCodeAlreadyExists
		}

		return Airport{}, fmt.Errorf("update airport: %w", err)
	}

	return *airport, nil
}