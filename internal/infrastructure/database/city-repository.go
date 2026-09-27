package database

import (
	"context"
	"errors"
	"fmt"

	"firefly-airline/internal/domains/city"

	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{
		db: db,
	}
}

var _ city.Repository = (*CityRepository)(nil)

func (r *CityRepository) Create(
	ctx context.Context,
	c *city.City,
) error {
	model := CityModel{
		ProvinceID: c.ProvinceID,
		Name:       c.Name,
	}

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return fmt.Errorf("create city: %w", err)
	}

	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	c.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *CityRepository) GetByID(
	ctx context.Context,
	id int64,
) (*city.City, error) {
	var model CityModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, city.ErrNotFound
		}

		return nil, fmt.Errorf("get city: %w", err)
	}

	return toDomainCity(&model), nil
}

func (r *CityRepository) List(
	ctx context.Context,
	provinceID *int64,
) ([]city.City, error) {
	var models []CityModel

	query := r.db.WithContext(ctx)

	if provinceID != nil {
		query = query.Where("province_id = ?", *provinceID)
	}

	if err := query.
		Order("id ASC").
		Find(&models).
		Error; err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}

	cities := make([]city.City, 0, len(models))

	for i := range models {
		cities = append(
			cities,
			*toDomainCity(&models[i]),
		)
	}

	return cities, nil
}

func (r *CityRepository) Update(
	ctx context.Context,
	c *city.City,
) error {
	var model CityModel

	err := r.db.WithContext(ctx).
		First(&model, c.ID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return city.ErrNotFound
		}

		return fmt.Errorf("find city: %w", err)
	}

	model.ProvinceID = c.ProvinceID
	model.Name = c.Name

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return fmt.Errorf("update city: %w", err)
	}

	c.UpdatedAt = model.UpdatedAt

	return nil
}

func toDomainCity(model *CityModel) *city.City {
	return &city.City{
		ID:         model.ID,
		ProvinceID: model.ProvinceID,
		Name:       model.Name,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}