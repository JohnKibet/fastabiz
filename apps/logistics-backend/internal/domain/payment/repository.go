package payment

import "github.com/google/uuid"

type Repository interface {
	Create(payment *Payment) error
	GetByID(id uuid.UUID) (*Payment, error)
	GetByOrder(id uuid.UUID) (*Payment, error)
	List() ([]*Payment, error)
}
