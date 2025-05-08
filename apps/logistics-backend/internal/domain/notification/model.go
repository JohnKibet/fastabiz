package notification

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	Email NotificationType = "email"
	SMS   NotificationType = "sms"
	Push  NotificationType = "push"
)

type Notification struct {
	ID      uuid.UUID        `db:"id" json:"id"`
	UserID  uuid.UUID        `db:"user_id" json:"user_id"`
	Message string           `db:"message" json:"message"`
	Type    NotificationType `db:"type" json:"type"`
	SentAt  time.Time        `db:"sent_at" json:"sent_at"`
}
