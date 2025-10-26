package user

import (
	"context"
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
)

type DriverReader interface {
	RegisterDriver(ctx context.Context, d *driver.Driver) error
}

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}
