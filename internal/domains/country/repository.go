package country

import "context"

type Repository interface {
	Create(ctx context.Context, country *Country) error
	GetByID(ctx context.Context, id int64) (*Country, error)
	List(ctx context.Context) ([]Country, error)
	Update(ctx context.Context, country *Country) error
}