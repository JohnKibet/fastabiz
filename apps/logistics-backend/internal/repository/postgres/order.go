package postgres

import (
	"logistics-backend/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) order.Repository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(o *order.Order) error {
	query := `
		INSERT INTO orders (customer_id, pickup_location, delivery_location, order_status)
		VALUES (:customer_id, :pickup_location, :delivery_location, :order_status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&o.ID, o)
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*order.Order, error) {
	query := `SELECT id FROM orders WHERE id = $1`
	var o order.Order
	err := r.db.Get(&o, query, id)
	return &o, err
}

func (r *OrderRepository) ListByCustomer(customerID uuid.UUID) ([]*order.Order, error) {
	query := `SELECT customer_id, pickup_location, delivery_location, order_status FROM orders WHERE customer_id = $1 ORDER BY created_at DESC`
	var orders []*order.Order
	err := r.db.Select(&orders, query, customerID)
	return orders, err
}

func (r *OrderRepository) UpdateStatus(orderID uuid.UUID, status order.OrderStatus) error {
	query := `UPDATE orders SET order_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, orderID)
	return err
}
