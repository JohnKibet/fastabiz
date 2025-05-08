package postgres

import (
	"logistics-backend/internal/domain/driver"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DriverRepository struct {
	db *sqlx.DB
}

func NewDriverRepository(db *sqlx.DB) driver.Repository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) Create(d *driver.Driver) error {
	query := `
		INSERT INTO drivers (full_name, email, vehicle_info, current_location, available)
		VALUES (:full_name, :email, :vehicle_info, :current_location, :available)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&d.ID, d)
}

func (r *DriverRepository) GetByID(id uuid.UUID) (*driver.Driver, error) {
	query := `SELECT id, full_name, email, vehicle_info, current_location, available FROM drivers WHERE id = $1`
	var d driver.Driver
	err := r.db.Get(&d, query, id)
	return &d, err
}

func (r *DriverRepository) GetByEmail(email string) (*driver.Driver, error) {
	query := `SELECT id, full_name, email, vehicle_info, current_location, available FROM drivers WHERE email = $1`
	var d driver.Driver
	err := r.db.Get(&d, query, email)
	return &d, err
}

func (r *DriverRepository) List() ([]*driver.Driver, error) {
	query := `SELECT id, full_name, available FROM drivers`
	var drivers []*driver.Driver
	err := r.db.Select(&drivers, query)
	return drivers, err
}
