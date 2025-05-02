package user

import "github.com/google/uuid"

type Repository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	List() ([]*User, error)
}
