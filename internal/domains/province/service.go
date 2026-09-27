package province

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"firefly-airline/internal/domains/country"
)

var ErrNotFound = errors.New("province not found")

type Service struct {
	repository        Repository
	countryRepository country.Repository
}

func NewService(
	repository Repository,
	countryRepository country.Repository,
) *Service {
	return &Service{
		repository:        repository,
		countryRepository: countryRepository,
	}
}

func (s *Service) Create(ctx context.Context, countryID int64, name string) (Province, error) {
	if countryID <= 0 {
		return Province{}, fmt.Errorf("invalid country id")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return Province{}, fmt.Errorf("province name is required")
	}

	if _, err := s.countryRepository.GetByID(ctx, countryID); err != nil {
		return Province{}, fmt.Errorf("country: %w", err)
	}

	province := &Province{
		CountryID: countryID,
		Name:      name,
	}

	if err := s.repository.Create(ctx, province); err != nil {
		return Province{}, fmt.Errorf("create province: %w", err)
	}

	return *province, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (Province, error) {
	province, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Province{}, ErrNotFound
		}

		return Province{}, fmt.Errorf("get province: %w", err)
	}

	return *province, nil
}

func (s *Service) List(ctx context.Context, countryID *int64) ([]Province, error) {
	if countryID != nil && *countryID <= 0 {
		return nil, fmt.Errorf("invalid country id")
	}

	provinces, err := s.repository.List(ctx, countryID)
	if err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}

	return provinces, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	countryID int64,
	name string,
) (Province, error) {
	if countryID <= 0 {
		return Province{}, fmt.Errorf("invalid country id")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return Province{}, fmt.Errorf("province name is required")
	}

	if _, err := s.countryRepository.GetByID(ctx, countryID); err != nil {
		return Province{}, fmt.Errorf("country: %w", err)
	}

	province := &Province{
		ID:        id,
		CountryID: countryID,
		Name:      name,
	}

	if err := s.repository.Update(ctx, province); err != nil {
		if errors.Is(err, ErrNotFound) {
			return Province{}, ErrNotFound
		}

		return Province{}, fmt.Errorf("update province: %w", err)
	}

	return *province, nil
}