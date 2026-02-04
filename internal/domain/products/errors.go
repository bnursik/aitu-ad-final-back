package products

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrNotFound           = errors.New("not found")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidCategory    = errors.New("invalid category")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrInvalidStock       = errors.New("invalid stock")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInvalidRating      = errors.New("invalid rating")
	ErrInvalidComment     = errors.New("invalid comment")
	ErrInvalidReviewID    = errors.New("invalid review id")
	ErrCannotDeleteProduct = errors.New("cannot delete product with stock; stock must be less than 1")
)
