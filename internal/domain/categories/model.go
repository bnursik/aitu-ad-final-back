package categories

import "time"

type Category struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListFilter struct {
	Offset int64
	Limit  int64
}

type CreateInput struct {
	Name        string
	Description string
}

type UpdateInput struct {
	Name        *string
	Description *string
}
