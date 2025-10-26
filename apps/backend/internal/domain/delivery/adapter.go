package delivery

import (
	"context"
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
	"backend/internal/domain/order"

	"github.com/google/uuid"
)

// cross-domain interface so delivery can access orders

type OrderReader interface {
	GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error)
	UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error
}

type DriverReader interface {
	GetDriverByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error)
	UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, value bool) error
}

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}
