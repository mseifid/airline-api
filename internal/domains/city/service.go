package city

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"firefly-airline/internal/domains/province"
)

var ErrNotFound = errors.New("city not found")

type Service struct {
	repository         Repository
	provinceRepository province.Repository
}

func NewService(
	repository Repository,
	provinceRepository province.Repository,
) *Service {
	return &Service{
		repository:         repository,
		provinceRepository: provinceRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	provinceID int64,
	name string,
) (City, error) {
	if provinceID <= 0 {
		return City{}, fmt.Errorf("invalid province id")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return City{}, fmt.Errorf("city name is required")
	}

	if _, err := s.provinceRepository.GetByID(ctx, provinceID); err != nil {
		return City{}, fmt.Errorf("province: %w", err)
	}

	city := &City{
		ProvinceID: provinceID,
		Name:       name,
	}

	if err := s.repository.Create(ctx, city); err != nil {
		return City{}, fmt.Errorf("create city: %w", err)
	}

	return *city, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (City, error) {
	city, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return City{}, ErrNotFound
		}

		return City{}, fmt.Errorf("get city: %w", err)
	}

	return *city, nil
}

func (s *Service) List(ctx context.Context, provinceID *int64) ([]City, error) {
	if provinceID != nil && *provinceID <= 0 {
		return nil, fmt.Errorf("invalid province id")
	}

	cities, err := s.repository.List(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	return cities, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	provinceID int64,
	name string,
) (City, error) {
	if provinceID <= 0 {
		return City{}, fmt.Errorf("invalid province id")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return City{}, fmt.Errorf("city name is required")
	}

	if _, err := s.provinceRepository.GetByID(ctx, provinceID); err != nil {
		return City{}, fmt.Errorf("province: %w", err)
	}

	city := &City{
		ID:         id,
		ProvinceID: provinceID,
		Name:       name,
	}

	if err := s.repository.Update(ctx, city); err != nil {
		if errors.Is(err, ErrNotFound) {
			return City{}, ErrNotFound
		}

		return City{}, fmt.Errorf("update city: %w", err)
	}

	return *city, nil
}