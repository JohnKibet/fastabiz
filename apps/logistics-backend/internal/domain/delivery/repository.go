package delivery

import "github.com/google/uuid"

type Repository interface {
	Create(deliver *Delivery) error
	GetByID(id uuid.UUID) (*Delivery, error)
	List() ([]*Delivery, error)
}
