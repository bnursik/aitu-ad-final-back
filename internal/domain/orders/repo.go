package orders

import "context"

type Repo interface {
	List(ctx context.Context, userID *string) ([]Order, error)
	GetByID(ctx context.Context, id string, userID *string) (Order, error)
	Create(ctx context.Context, o Order) (Order, error)
	UpdateStatus(ctx context.Context, id string, status Status) (Order, error)
}
