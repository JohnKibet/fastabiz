package driver

import "github.com/google/uuid"

type Repository interface {
	Create(driver *Driver) error
	GetByID(id uuid.UUID) (*Driver, error)
	GetByEmail(email string) (*Driver, error)
	List() ([]*Driver, error)
}
