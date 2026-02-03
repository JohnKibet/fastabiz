package driver

import "errors"

var (
	ErrMissingUserID       = errors.New("missing driver ID")
	ErrDriverAlreadyExists = errors.New("An account with this email already exists.")
	ErrRoleCheck              = errors.New("Invalid user role.")
)
