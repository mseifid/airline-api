package airplane

import "context"

type Repository interface {
	Create(ctx context.Context, airplane *Airplane) error
	GetByID(ctx context.Context, id int64) (*Airplane, error)
	List(ctx context.Context) ([]Airplane, error)
	Update(ctx context.Context, airplane *Airplane) error
	Delete(ctx context.Context, id int64) error
}