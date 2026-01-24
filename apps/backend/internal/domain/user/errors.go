package user

import "errors"

var (
	ErrInvalidPhone           = errors.New("Invalid phone number.")
	ErrInvalidEmail           = errors.New("Invalid email address.")
	ErrInvalidName            = errors.New("Invalid name.")
	ErrDB                     = errors.New("Database error.")
	ErrCreateUser             = errors.New("Register user failed.")
	ErrUserAlreadyExists      = errors.New("An account with this email already exists.")
	ErrRoleCheck              = errors.New("Invalid user role.")
	ErrInvalidUserId          = errors.New("Invalid user id.")
	ErrInvalidCurrentPassword = errors.New("Current Password is incorrect.")
	ErrInvalidStatusInput     = errors.New("Invalid user status input.")
	ErrInvalidDataInput       = errors.New("Invalid data input.")
	ErrInvalidColumn          = errors.New("Invalid column update.")
	ErrUserHasReferences      = errors.New("User has dependent records.")
	ErrUserNotFound           = errors.New("User not found.")
)
