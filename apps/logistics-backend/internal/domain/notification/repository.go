package notification

import "github.com/google/uuid"

type Repository interface {
	Create(notification *Notification) error
	GetByID(id uuid.UUID) (*Notification, error)
	List() ([]*Notification, error)
}
