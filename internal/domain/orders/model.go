package orders

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

type Order struct {
	ID        string
	UserID    string
	Items     []Item
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ProductID string
	Quantity  int64
}

type CreateInput struct {
	Items []Item
}
