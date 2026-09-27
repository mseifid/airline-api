package city

import "context"

type Repository interface {
	Create(ctx context.Context, city *City) error
	GetByID(ctx context.Context, id int64) (*City, error)
	List(ctx context.Context, provinceID *int64) ([]City, error)
	Update(ctx context.Context, city *City) error
}