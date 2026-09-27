package flight

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, flight *Flight) error
	GetByID(ctx context.Context, id int64,) (*Flight, error)
	Update(ctx context.Context, flight *Flight) error
	ListActiveByCitiesAndDateRange(
		ctx context.Context,
		sourceCityID int64,
		destinationCityID int64,
		from time.Time,
		to time.Time,
	) ([]Flight, error)
	ListActive(ctx context.Context) ([]Flight, error)
}