package airport

import "context"

type Repository interface {
	Create(ctx context.Context, airport *Airport) error
	GetByID(ctx context.Context, id int64) (*Airport, error)
	List(ctx context.Context, cityID *int64) ([]Airport, error)
	Update(ctx context.Context, airport *Airport) error
}