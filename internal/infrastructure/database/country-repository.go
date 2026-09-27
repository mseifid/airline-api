package database

import (
	"context"
	"errors"
	"fmt"

	"firefly-airline/internal/domains/country"

	"gorm.io/gorm"
)

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{
		db: db,
	}
}

var _ country.Repository = (*CountryRepository)(nil)

func (r *CountryRepository) Create(ctx context.Context, c *country.Country) error {
	model := CountryModel{
		Name: c.Name,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create country: %w", err)
	}

	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	c.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *CountryRepository) GetByID(ctx context.Context, id int64) (*country.Country, error) {
	var model CountryModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, country.ErrNotFound
		}

		return nil, fmt.Errorf("get country: %w", err)
	}

	return toDomainCountry(&model), nil
}

func (r *CountryRepository) List(ctx context.Context) ([]country.Country, error) {
	var models []CountryModel

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&models).
		Error

	if err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}

	countries := make([]country.Country, 0, len(models))

	for i := range models {
		countries = append(
			countries,
			*toDomainCountry(&models[i]),
		)
	}

	return countries, nil
}

func (r *CountryRepository) Update(ctx context.Context, c *country.Country) error {
	var model CountryModel

	err := r.db.WithContext(ctx).
		First(&model, c.ID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return country.ErrNotFound
		}

		return fmt.Errorf("find country: %w", err)
	}

	model.Name = c.Name

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("update country: %w", err)
	}

	c.UpdatedAt = model.UpdatedAt

	return nil
}

func toDomainCountry(model *CountryModel) *country.Country {
	return &country.Country{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}