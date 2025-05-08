package postgres

import (
	"logistics-backend/internal/domain/feedback"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FeedbackRepository struct {
	db *sqlx.DB
}

func NewFeedbackRepository(db *sqlx.DB) feedback.Repository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(f *feedback.Feedback) error {
	query := `
		INSERT INTO feedbacks (order_id, customer_id, rating, comments)
		VALUES (:order_id, :customer_id, :rating, :comments)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&f.ID, f)
}

func (r *FeedbackRepository) GetByID(id uuid.UUID) (*feedback.Feedback, error) {
	query := `SELECT id, order_id, customer_id, rating, comments FROM feedbacks WHERE id = $1`
	var f feedback.Feedback
	err := r.db.Get(&f, query, id)
	return &f, err
}

func (r *FeedbackRepository) List() ([]*feedback.Feedback, error) {
	query := `SELECT id, order_id, customer_id, rating, comments FROM feedbacks`
	var feedbacks []*feedback.Feedback
	err := r.db.Select(&feedbacks, query)
	return feedbacks, err
}
