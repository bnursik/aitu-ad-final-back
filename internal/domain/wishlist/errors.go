package wishlist

import "errors"

var (
	ErrInvalidID      = errors.New("invalid id")
	ErrInvalidProduct = errors.New("invalid product")
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
)
