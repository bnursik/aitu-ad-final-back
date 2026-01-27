package categories

import "context"

type Repo interface {
	List(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id string) (Category, error)
	Create(ctx context.Context, c Category) (Category, error)
	Update(ctx context.Context, id string, in UpdateInput) (Category, error)
	Delete(ctx context.Context, id string) error
}
