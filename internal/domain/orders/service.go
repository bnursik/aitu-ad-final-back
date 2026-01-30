package orders

import "context"

type Service interface {
	List(ctx context.Context, userID string, isAdmin bool, f ListFilter) ([]Order, error)
	Get(ctx context.Context, id string, userID string, isAdmin bool) (Order, error)
	Create(ctx context.Context, userID string, in CreateInput) (Order, error)
	UpdateStatus(ctx context.Context, id string, status Status) (Order, error) // admin only (проверяй в handler)
}
