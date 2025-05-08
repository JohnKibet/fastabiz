package postgres

import (
	"fmt"
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
		INSERT INTO orders (user_id, pickup_address, delivery_address, status)
		VALUES (:user_id, :pickup_address, :delivery_address, :status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&o.ID, o)
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*order.Order, error) {
	query := `SELECT id, user_id, status FROM orders WHERE id = $1`
	var o order.Order
	err := r.db.Get(&o, query, id)
	return &o, err
}

func (r *OrderRepository) ListByCustomer(customerID uuid.UUID) ([]*order.Order, error) {
	query := `SELECT id, user_id, status FROM orders WHERE user_id = $1`
	var orders []*order.Order
	err := r.db.Select(&orders, query, customerID)
	return orders, err
}

func (r *OrderRepository) UpdateStatus(orderID uuid.UUID, status order.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order found with id %s", orderID)
	}
	return nil
}

func (r *OrderRepository) List() ([]*order.Order, error) {
	query := `SELECT id, user_id, pickup_address, delivery_address, status FROM orders`
	var orders []*order.Order
	err := r.db.Select(&orders, query)
	return orders, err
}
