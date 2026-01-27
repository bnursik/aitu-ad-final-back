package categories

import "context"

type ProductsCounter interface {
	CountByCategoryID(ctx context.Context, categoryID string) (int64, error)
}

type Service interface {
	List(ctx context.Context) ([]Category, error)
	Get(ctx context.Context, id string) (Category, error)
	Create(ctx context.Context, in CreateInput) (Category, error)
	Update(ctx context.Context, id string, in UpdateInput) (Category, error)
	Delete(ctx context.Context, id string) error
}
