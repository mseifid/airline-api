package database

import (
	"context"
	"errors"
	"fmt"

	"firefly-airline/internal/domains/province"

	"gorm.io/gorm"
)

type ProvinceRepository struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) *ProvinceRepository {
	return &ProvinceRepository{
		db: db,
	}
}

var _ province.Repository = (*ProvinceRepository)(nil)

func (r *ProvinceRepository) Create(
	ctx context.Context,
	p *province.Province,
) error {
	model := ProvinceModel{
		CountryID: p.CountryID,
		Name:      p.Name,
	}

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		return fmt.Errorf("create province: %w", err)
	}

	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *ProvinceRepository) GetByID(
	ctx context.Context,
	id int64,
) (*province.Province, error) {
	var model ProvinceModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, province.ErrNotFound
		}

		return nil, fmt.Errorf("get province: %w", err)
	}

	return toDomainProvince(&model), nil
}

func (r *ProvinceRepository) List(
	ctx context.Context,
	countryID *int64,
) ([]province.Province, error) {
	var models []ProvinceModel

	query := r.db.WithContext(ctx)

	if countryID != nil {
		query = query.Where("country_id = ?", *countryID)
	}

	if err := query.
		Order("id ASC").
		Find(&models).
		Error; err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}

	provinces := make([]province.Province, 0, len(models))

	for i := range models {
		provinces = append(
			provinces,
			*toDomainProvince(&models[i]),
		)
	}

	return provinces, nil
}

func (r *ProvinceRepository) Update(
	ctx context.Context,
	p *province.Province,
) error {
	var model ProvinceModel

	err := r.db.WithContext(ctx).
		First(&model, p.ID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return province.ErrNotFound
		}

		return fmt.Errorf("find province: %w", err)
	}

	model.CountryID = p.CountryID
	model.Name = p.Name

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		return fmt.Errorf("update province: %w", err)
	}

	p.UpdatedAt = model.UpdatedAt

	return nil
}

func toDomainProvince(model *ProvinceModel) *province.Province {
	return &province.Province{
		ID:        model.ID,
		CountryID: model.CountryID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}