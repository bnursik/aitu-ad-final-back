package orders

import "errors"

var (
	ErrInvalidID      = errors.New("invalid id")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidItems   = errors.New("invalid items")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidProduct = errors.New("invalid product")
	ErrInvalidQty     = errors.New("invalid quantity")
)
