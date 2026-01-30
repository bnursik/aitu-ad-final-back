package wishlist

import "time"

type WishlistItem struct {
	ID        string
	UserID    string
	ProductID string
	CreatedAt time.Time
}

type AddItemInput struct {
	ProductID string
}

type ListFilter struct {
	Offset int64
	Limit  int64
}
