package products

import "time"

type Product struct {
	ID          string
	CategoryID  string
	Name        string
	Description string
	Price       float64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Reviews     []Review
}

type Review struct {
	ID        string
	UserID    string
	Rating    int64
	Comment   string
	CreatedAt time.Time
}

type ListFilter struct {
	CategoryID *string
}

type CreateInput struct {
	CategoryID  string
	Name        string
	Description string
	Price       float64
	Stock       int64
}

type UpdateInput struct {
	CategoryID  *string
	Name        *string
	Description *string
	Price       *float64
	Stock       *int64
}

type AddReviewInput struct {
	UserID  string
	Rating  int64
	Comment string
}
