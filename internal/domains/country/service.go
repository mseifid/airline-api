package country

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("country not found")

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
	name string,
) (Country, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Country{}, fmt.Errorf("country name is required")
	}

	country := &Country{
		Name: name,
	}

	if err := s.repository.Create(ctx, country); err != nil {
		return Country{}, fmt.Errorf("create country: %w", err)
	}

	return *country, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (Country, error) {
	country, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Country{}, ErrNotFound
		}

		return Country{}, fmt.Errorf("get country: %w", err)
	}

	return *country, nil
}

func (s *Service) List(
	ctx context.Context,
) ([]Country, error) {
	countries, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}

	return countries, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	name string,
) (Country, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Country{}, fmt.Errorf("country name is required")
	}

	country := &Country{
		ID:   id,
		Name: name,
	}

	if err := s.repository.Update(ctx, country); err != nil {
		if errors.Is(err, ErrNotFound) {
			return Country{}, ErrNotFound
		}

		return Country{}, fmt.Errorf("update country: %w", err)
	}

	return *country, nil
}