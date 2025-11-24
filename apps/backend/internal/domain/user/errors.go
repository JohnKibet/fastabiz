package user

import "errors"

var (
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrInvalidEmail = errors.New("invalid email address")
	ErrInvalidName  = errors.New("invalid name")
	ErrDB           = errors.New("database error")
)
