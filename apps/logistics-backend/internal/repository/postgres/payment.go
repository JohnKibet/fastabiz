package postgres

import (
	"logistics-backend/internal/domain/payment"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) payment.Repository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(p *payment.Payment) error {
	query := `
		INSERT INTO payments (order_id, amount, method, status)
		VALUES (:order_id, :amount, :method, :status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.Get(&p.ID, p)
}

func (r *PaymentRepository) GetByID(id uuid.UUID) (*payment.Payment, error) {
	query := `SELECT id, order_id, amount, method, status, paid_at FROM payments WHERE id = $1`
	var p payment.Payment
	err := r.db.Get(&p, query, id)
	return &p, err
}

func (r *PaymentRepository) GetByOrder(order_id uuid.UUID) (*payment.Payment, error) {
	query := `SELECT id, order_id, amount, status, paid_at FROM payments WHERE order_id = $1` // get customer_id via order_id
	var p payment.Payment
	err := r.db.Get(&p, query, order_id)
	return &p, err
}

func (r *PaymentRepository) List() ([]*payment.Payment, error) {
	query := `SELECT id, order_id, amount, status FROM payments`
	var payments []*payment.Payment
	err := r.db.Select(&payments, query)
	return payments, err
}
