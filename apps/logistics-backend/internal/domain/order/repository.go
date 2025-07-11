package order

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(order *Order) error                                                           //POST method to create order.
	GetByID(id uuid.UUID) (*Order, error)                                                //GET method for fetching order by id.
	ListByCustomer(customerID uuid.UUID) ([]*Order, error)                               //GET method for fetching all orders by customer id.
	UpdateColumn(ctx context.Context, orderID uuid.UUID, column string, value any) error // PUT generic method to update specified column value in orders table.
	List() ([]*Order, error)                                                             // GET method for fetching all orders

	// UpdateStatus(orderID uuid.UUID, status OrderStatus) error //PATCH method to update order status.
	// PUT
	//DELETE
}
