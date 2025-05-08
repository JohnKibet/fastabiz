package postgres

import (
	"logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) delivery.Repository {
	return &DeliveryRepository{db: db}
}

func (r *DeliveryRepository) Create(d *delivery.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, driver_id, status)
		VALUES (:order_id, :driver_id, :status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&d.ID, d)
}

func (r *DeliveryRepository) GetByID(id uuid.UUID) (*delivery.Delivery, error) {
	query := `SELECT id, order_id, driver_id, status FROM deliveries WHERE id = $1`
	var d delivery.Delivery
	err := r.db.Get(&d, query, id)
	return &d, err
}

func (r *DeliveryRepository) List() ([]*delivery.Delivery, error) {
	query := `SELECT id, order_id, driver_id, status FROM deliveries`
	var deliveries []*delivery.Delivery
	err := r.db.Select(&deliveries, query)
	return deliveries, err
}
