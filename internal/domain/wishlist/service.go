package wishlist

import "context"

type Service interface {
	Add(ctx context.Context, userID string, in AddItemInput) (WishlistItem, error)
	List(ctx context.Context, userID string, f ListFilter) ([]WishlistItem, int64, error)
	Delete(ctx context.Context, userID string, id string) error
}
