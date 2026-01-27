package users

import "errors"

var (
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)
