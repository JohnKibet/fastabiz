package user

import "errors"

var (
	ErrInvalidPhone      = errors.New("Invalid phone number.")
	ErrInvalidEmail      = errors.New("Invalid email address.")
	ErrInvalidName       = errors.New("Invalid name.")
	ErrDB                = errors.New("Database error.")
	ErrCreateUser        = errors.New("Register user failed.")
	ErrUserAlreadyExists = errors.New("An account with this email already exists.")
)
