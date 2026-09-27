package database

import (
	"context"
	"errors"
	"fmt"

	"firefly-airline/internal/domains/airport"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AirportRepository struct {
	db *gorm.DB
}

func NewAirportRepository(db *gorm.DB) *AirportRepository {
	return &AirportRepository{
		db: db,
	}
}

var _ airport.Repository = (*AirportRepository)(nil)

func (r *AirportRepository) Create(
	ctx context.Context,
	a *airport.Airport,
) error {
	model := AirportModel{
		CityID: a.CityID,
		Code:   a.Code,
		Name:   a.Name,
	}

	if err := r.db.WithContext(ctx).
		Create(&model).
		Error; err != nil {
		if isUniqueViolation(err) {
			return airport.ErrCodeAlreadyExists
		}

		return fmt.Errorf("create airport: %w", err)
	}

	a.ID = model.ID
	a.CreatedAt = model.CreatedAt
	a.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *AirportRepository) GetByID(
	ctx context.Context,
	id int64,
) (*airport.Airport, error) {
	var model AirportModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, airport.ErrNotFound
		}

		return nil, fmt.Errorf("get airport: %w", err)
	}

	return toDomainAirport(&model), nil
}

func (r *AirportRepository) List(
	ctx context.Context,
	cityID *int64,
) ([]airport.Airport, error) {
	var models []AirportModel

	query := r.db.WithContext(ctx)

	if cityID != nil {
		query = query.Where(
			"city_id = ?",
			*cityID,
		)
	}

	if err := query.
		Order("id ASC").
		Find(&models).
		Error; err != nil {
		return nil, fmt.Errorf("list airports: %w", err)
	}

	airports := make(
		[]airport.Airport,
		0,
		len(models),
	)

	for i := range models {
		airports = append(
			airports,
			*toDomainAirport(&models[i]),
		)
	}

	return airports, nil
}

func (r *AirportRepository) Update(
	ctx context.Context,
	a *airport.Airport,
) error {
	var model AirportModel

	err := r.db.WithContext(ctx).
		First(&model, a.ID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return airport.ErrNotFound
		}

		return fmt.Errorf("find airport: %w", err)
	}

	model.CityID = a.CityID
	model.Code = a.Code
	model.Name = a.Name

	if err := r.db.WithContext(ctx).
		Save(&model).
		Error; err != nil {
		if isUniqueViolation(err) {
			return airport.ErrCodeAlreadyExists
		}

		return fmt.Errorf("update airport: %w", err)
	}

	a.UpdatedAt = model.UpdatedAt

	return nil
}

func toDomainAirport(
	model *AirportModel,
) *airport.Airport {
	return &airport.Airport{
		ID:        model.ID,
		CityID:    model.CityID,
		Code:      model.Code,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) &&
		pgErr.Code == "23505"
}