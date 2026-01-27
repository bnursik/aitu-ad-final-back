package products

import "context"

type Service interface {
	List(ctx context.Context, f ListFilter) ([]Product, error)
	Get(ctx context.Context, id string) (Product, error)
	Create(ctx context.Context, in CreateInput) (Product, error)
	Update(ctx context.Context, id string, in UpdateInput) (Product, error)
	Delete(ctx context.Context, id string) error

	AddReview(ctx context.Context, productID string, in AddReviewInput) (Review, error)
	DeleteReview(ctx context.Context, productID, reviewID string) error
}
