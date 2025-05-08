package notification

import "github.com/google/uuid"

type CreateNotificationRequest struct {
	UserID  uuid.UUID        `json:"user_id"`
	Message string           `json:"message"`
	Type    NotificationType `json:"type"`
}

func (r *CreateNotificationRequest) ToNotification() *Notification {
	return &Notification{
		UserID:  r.UserID,
		Message: r.Message,
		Type:    r.Type,
	}
}
