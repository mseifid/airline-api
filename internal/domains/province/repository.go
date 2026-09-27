package province

import "context"

type Repository interface {
	Create(ctx context.Context, province *Province) error
	GetByID(ctx context.Context, id int64) (*Province, error)
	List(ctx context.Context, countryID *int64) ([]Province, error)
	Update(ctx context.Context, province *Province) error
}