package airplane

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, airplaneType string, capacity int) (*Airplane, error) {
	airplaneType = strings.TrimSpace(airplaneType)

	if airplaneType == "" {
		return nil, fmt.Errorf("type is required")
	}

	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be greater than zero")
	}

	airplane := &Airplane{
		Type:     airplaneType,
		Capacity: capacity,
	}

	if err := s.repository.Create(ctx, airplane); err != nil {
		return nil, fmt.Errorf("create airplane: %w", err)
	}

	return airplane, nil
}

func (s *Service) GetByID(ctx context.Context,id int64) (*Airplane, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid airplane id")
	}

	airplane, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get airplane: %w", err)
	}

	return airplane, nil
}

func (s *Service) List(ctx context.Context,) ([]Airplane, error) {
	airplanes, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list airplanes: %w", err)
	}

	return airplanes, nil
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	airplaneType string,
	capacity int,
) (*Airplane, error) {
	airplaneType = strings.TrimSpace(airplaneType)

	if id <= 0 {
		return nil, fmt.Errorf("invalid airplane id")
	}

	if airplaneType == "" {
		return nil, fmt.Errorf("type is required")
	}

	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be greater than zero")
	}

	airplane := &Airplane{
		ID:       id,
		Type:     airplaneType,
		Capacity: capacity,
	}

	if err := s.repository.Update(ctx, airplane); err != nil {
		return nil, fmt.Errorf("update airplane: %w", err)
	}

	return airplane, nil
}

func (s *Service) Delete(ctx context.Context,id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid airplane id")
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete airplane: %w", err)
	}

	return nil
}
