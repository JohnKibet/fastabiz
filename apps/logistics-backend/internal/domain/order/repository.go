package order

import "github.com/google/uuid"

type Repository interface {
	Create(order *Order) error
	GetByID(id uuid.UUID) (*Order, error)
	ListByCustomer(customerID uuid.UUID) ([]*Order, error)
	UpdateStatus(orderID uuid.UUID, status OrderStatus) error
	List() ([]*Order, error)
}
