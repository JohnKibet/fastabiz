package order

import (
	"context"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

// Repository defines CRUD and lookup operations for orders.
type Repository interface {
	// Create inserts a new order record
	Create(ctx context.Context, order *Order) error

	// GetByID fetches a single order by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)

	// ListByCustomer returns all orders for a given customer
	ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*Order, error)

	// Update updates a single column for a given order
	Update(ctx context.Context, orderID uuid.UUID, column string, value any) error

	// List returns all orders
	List(ctx context.Context) ([]*Order, error)

	// Delete removes an order by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetPickupPoint returns the pickup location of an order
	GetPickupPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error)

	// GetDeliveryPoint returns the delivery location of an order
	GetDeliveryPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error)

	// CreatePending inserts a pending order
	CreatePending(ctx context.Context, o *Order) error
}
