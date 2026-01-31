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
	ID         string
	UserID     string
	Items      []Item
	Status     Status
	TotalPrice float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ProductID string
	Quantity  int64

	UnitPrice float64
	LineTotal float64
}

type CreateInput struct {
	Items []Item
}

type ListFilter struct {
	Offset int64
	Limit  int64
}
