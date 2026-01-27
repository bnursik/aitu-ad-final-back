package categories

import "errors"

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrInvalidName = errors.New("invalid name")
	ErrNotFound    = errors.New("not found")
	ErrHasProducts = errors.New("category has products")
)
