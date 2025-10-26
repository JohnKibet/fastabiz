package inventory

import (
	"context"
	"backend/internal/domain/notification"
	"backend/internal/domain/store"

	"github.com/google/uuid"
)

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}

type StoreReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*store.Store, error)
}
