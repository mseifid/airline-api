package database

import (
	"context"
	"errors"
	"fmt"

	"firefly-airline/internal/domains/airplane"

	"gorm.io/gorm"
)

type AirplaneRepository struct {
	db *gorm.DB
}

func NewAirplaneRepository(db *gorm.DB) *AirplaneRepository {
	return &AirplaneRepository{
		db: db,
	}
}

var _ airplane.Repository = (*AirplaneRepository)(nil)

func (r *AirplaneRepository) Create(ctx context.Context, a *airplane.Airplane) error {
	model := AirplaneModel{
		Type:     a.Type,
		Capacity: int16(a.Capacity),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create airplane: %w", err)
	}

	a.ID = int64(model.ID)
	a.CreatedAt = model.CreatedAt
	a.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *AirplaneRepository) GetByID(ctx context.Context, id int64) (*airplane.Airplane, error) {
	var model AirplaneModel

	if err := r.db.WithContext(ctx).
		First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("airplane not found")
		}

		return nil, fmt.Errorf("get airplane: %w", err)
	}

	return toAirplaneDomain(model), nil
}

func (r *AirplaneRepository) List(ctx context.Context) ([]airplane.Airplane, error) {
	var models []AirplaneModel

	if err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list airplanes: %w", err)
	}

	airplanes := make([]airplane.Airplane, 0, len(models))

	for _, model := range models {
		airplanes = append(airplanes, *toAirplaneDomain(model))
	}

	return airplanes, nil
}

func (r *AirplaneRepository) Update(ctx context.Context, a *airplane.Airplane) error {
	result := r.db.WithContext(ctx).
		Model(&AirplaneModel{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"type":     a.Type,
			"capacity": a.Capacity,
		})

	if result.Error != nil {
		return fmt.Errorf("update airplane: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("airplane not found")
	}

	return nil
}

func (r *AirplaneRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Delete(&AirplaneModel{}, id)

	if result.Error != nil {
		return fmt.Errorf("delete airplane: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("airplane not found")
	}

	return nil
}

func toAirplaneDomain(model AirplaneModel) *airplane.Airplane {
	return &airplane.Airplane{
		ID:        int64(model.ID),
		Type:      model.Type,
		Capacity:  int(model.Capacity),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}