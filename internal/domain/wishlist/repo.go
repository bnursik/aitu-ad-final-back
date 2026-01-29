package wishlist

import "context"

type Repo interface {
	Add(ctx context.Context, item WishlistItem) (WishlistItem, error)
	List(ctx context.Context, userID string) ([]WishlistItem, error)
	Delete(ctx context.Context, userID string, id string) error
	Exists(ctx context.Context, userID string, productID string) (bool, error)
}
