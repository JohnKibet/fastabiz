package postgres

import (
	"logistics-backend/internal/domain/notification"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) notification.Repository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *notification.Notification) error {
	query := `
		INSERT INTO notifications (user_id, message, type)
		VALUES (:user_id, :message, :type)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.Get(&n.ID, n)
}

func (r *NotificationRepository) GetByID(id uuid.UUID) (*notification.Notification, error) {
	query := `SELECT id, user_id, message, type FROM notifications WHERE id = $1`
	var n notification.Notification
	err := r.db.Get(&n, query, id)
	return &n, err
}

func (r *NotificationRepository) List() ([]*notification.Notification, error) {
	query := `SELECT id, user_id, message, type FROM notifications`
	var notifications []*notification.Notification
	err := r.db.Select(&notifications, query)
	return notifications, err
}
