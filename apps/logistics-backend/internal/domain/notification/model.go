package notification

type NotificationType string

const (
	Email NotificationType = "email"
	SMS   NotificationType = "sms"
	Push  NotificationType = "push"
)

type Notification struct {
	ID      int64            `json:"id"`
	UserID  int64            `json:"user_id"`
	Message string           `json:"message"`
	Type    NotificationType `json:"type"`
	SentAt  string           `json:"sent_at"`
}
