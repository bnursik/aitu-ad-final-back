package products

import "errors"

var (
	ErrInvalidID       = errors.New("invalid id")
	ErrNotFound        = errors.New("not found")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidCategory = errors.New("invalid category")
	ErrInvalidPrice    = errors.New("invalid price")
	ErrInvalidStock    = errors.New("invalid stock")
	ErrInvalidRating   = errors.New("invalid rating")
	ErrInvalidComment  = errors.New("invalid comment")
	ErrInvalidReviewID = errors.New("invalid review id")
)
