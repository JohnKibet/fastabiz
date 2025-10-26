package postgres

import (
	"context"
	"fmt"
	"backend/internal/application"
	"backend/internal/domain/order"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	exec sqlx.ExtContext
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{exec: db}
}

func (r *OrderRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	query := `
		INSERT INTO orders (
		admin_id, user_id, inventory_id, quantity,
		pickup_address, pickup_point,
		delivery_address, delivery_point,
		status
	) VALUES (
		:admin_id, :user_id, :inventory_id, :quantity,
		:pickup_address, ST_SetSRID(ST_MakePoint(:pickup_point.x, :pickup_point.y), 4326),
		:delivery_address, ST_SetSRID(ST_MakePoint(:delivery_point.x, :delivery_point.y), 4326),
		:status
	)
	RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, o)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&o.ID); err != nil {
			return fmt.Errorf("scanning new order id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	query := `
		SELECT id, user_id, admin_id, inventory_id, quantity, pickup_address, pickup_point, delivery_address, delivery_point, status, created_at, updated_at 
		FROM orders 
		WHERE id = $1
	`

	var o order.Order
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &o, query, id); err != nil {
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	return &o, nil
}

func (r *OrderRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	query := `
		SELECT id, user_id, admin_id, inventory_id, quantity, pickup_address, pickup_point, delivery_address, delivery_point, status, created_at, updated_at 
		FROM orders 
		WHERE user_id = $1
	`

	var orders []*order.Order
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &orders, query, customerID)
	return orders, err
}

func (r *OrderRepository) Update(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	allowed := map[string]bool{
		"status":           true,
		"quantity":         true,
		"pickup_address":   true,
		"delivery_address": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE orders SET %s = :value, updated_at = NOW() 
		WHERE id = :id
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    orderID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", orderID)
	}

	return nil
}

// func (r *OrderRepository) UpdateStatus(orderID uuid.UUID, status order.OrderStatus) error {
// 	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
// 	res, err := r.db.Exec(query, status, orderID)
// 	if err != nil {
// 		return err
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("no order found with id %s", orderID)
// 	}
// 	return nil
// }

func (r *OrderRepository) List(ctx context.Context) ([]*order.Order, error) {
	query := `
		SELECT id, user_id, admin_id, inventory_id, quantity, pickup_address, delivery_address, status, created_at, updated_at, pickup_point, delivery_point
		FROM orders
	`

	var orders []*order.Order
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &orders, query)
	return orders, err
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM orders 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify order deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("order already deleted or invalid")
	}

	return nil
}

func (r *OrderRepository) GetPickupPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	var pt postgis.PointS
	query := `
        SELECT pickup_point
        FROM orders
        WHERE id = $1
    `
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &pt, query, orderID)
	if err != nil {
		return postgis.PointS{}, fmt.Errorf("get pickup point: %w", err)
	}
	return pt, nil
}

func (r *OrderRepository) GetDeliveryPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	var pt postgis.PointS
	query := `
        SELECT delivery_point
        FROM orders
        WHERE id = $1
    `
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &pt, query, orderID)
	if err != nil {
		return postgis.PointS{}, fmt.Errorf("get delivery point: %w", err)
	}
	return pt, nil
}
