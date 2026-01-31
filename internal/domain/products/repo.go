package products

import "context"

type Repo interface {
	List(ctx context.Context, f ListFilter) ([]Product, error)
	Count(ctx context.Context, f ListFilter) (int64, error)
	GetByID(ctx context.Context, id string) (Product, error)
	Create(ctx context.Context, p Product) (Product, error)
	Update(ctx context.Context, id string, in UpdateInput) (Product, error)
	Delete(ctx context.Context, id string) error

	AddReview(ctx context.Context, productID string, r Review) (Review, error)
	DeleteReview(ctx context.Context, productID string, reviewID string) error
}
