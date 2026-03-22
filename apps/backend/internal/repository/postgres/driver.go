package postgres

import (
	"backend/internal/application"
	"backend/internal/domain/driver"
	"context"
	"fmt"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DriverRepository struct {
	exec sqlx.ExtContext
}

func NewDriverRepository(db *sqlx.DB) *DriverRepository {
	return &DriverRepository{exec: db}
}

func (r *DriverRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *DriverRepository) Create(ctx context.Context, d *driver.Driver) error {
	//
	// FIX: Do NOT pass the raw *driver.Driver struct to sqlx.NamedQueryContext when it
	// contains a postgis.PointS field.
	//
	// Root cause of "invalid byte sequence for encoding UTF8: 0x00":
	//   sqlx.NamedQueryContext binds struct fields by their db tags. When it encounters
	//   postgis.PointS it calls the Value() method which produces an EWKB binary blob.
	//   EWKB is a raw binary format — it contains 0x00 null bytes by design.
	//   Postgres then tries to store that binary string in a geometry column through
	//   the text protocol, hits the null byte, and rejects it with the UTF-8 error.
	//
	// Fix: build the WKT string manually and pass a plain map[string]interface{}.
	//   WKT (Well-Known Text) is a human-readable ASCII format — no null bytes.
	//   ST_GeomFromEWKT() on the SQL side parses it back into a proper geometry.
	//
	// This is the same pattern already used in UpdateProfile — it was just never
	// applied to Create.
	//
	wkt := fmt.Sprintf(
		"SRID=%d;POINT(%f %f)",
		d.CurrentLocation.SRID,
		d.CurrentLocation.X, // longitude
		d.CurrentLocation.Y, // latitude
	)

	query := `
		INSERT INTO drivers (id, full_name, email, vehicle_info, current_location, available, created_at)
		VALUES (:id, :full_name, :email, :vehicle_info, ST_GeomFromEWKT(:current_location), :available, :created_at)
		RETURNING id
	`

	args := map[string]interface{}{
		"id":               d.ID,
		"full_name":        d.FullName,
		"email":            d.Email,
		"vehicle_info":     d.VehicleInfo,
		"current_location": wkt, // plain ASCII string — no binary, no null bytes
		"available":        d.Available,
		"created_at":       d.CreatedAt,
	}

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return driver.ErrDriverAlreadyExists
			case "23514":
				return driver.ErrRoleCheck
			case "22021":
				return fmt.Errorf("invalid characters in input: %w", err)
			}
		}
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&d.ID); err != nil {
			return fmt.Errorf("scanning new driver id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after insert")
	}

	return nil
}

// UpdateProfile — already uses WKT correctly, no changes needed.
func (r *DriverRepository) UpdateProfile(ctx context.Context, driverID uuid.UUID, vehicleInfo string, currentLocation postgis.PointS) error {
	query := `
		UPDATE drivers 
		SET vehicle_info = :vehicle, current_location = ST_GeomFromEWKT(:location)
		WHERE id = :id
	`

	wkt := fmt.Sprintf("SRID=%d;POINT(%f %f)", currentLocation.SRID, currentLocation.X, currentLocation.Y)

	args := map[string]interface{}{
		"vehicle":  vehicleInfo,
		"location": wkt,
		"id":       driverID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update driver profile: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no driver found with id %s", driverID)
	}
	return nil
}

func (r *DriverRepository) UpdateColumn(ctx context.Context, driverID uuid.UUID, column string, value any) error {
	allowed := map[string]bool{
		"full_name":        true,
		"email":            true,
		"vehicle_info":     true,
		"current_location": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE drivers SET %s = $1, updated_at = NOW() 
		WHERE id = $2
	`, column)

	res, err := r.execFromCtx(ctx).ExecContext(ctx, query, value, driverID)
	if err != nil {
		return fmt.Errorf("update driver: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no driver found with id %s", driverID)
	}
	return nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers 
		WHERE id = $1
	`
	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, id)
	return &d, err
}

func (r *DriverRepository) GetByEmail(ctx context.Context, email string) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers 
		WHERE email = $1
	`
	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, email)
	return &d, err
}

func (r *DriverRepository) List(ctx context.Context) ([]*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers
	`
	var drivers []*driver.Driver
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &drivers, query)
	return drivers, err
}

func (r *DriverRepository) ListAvailableDrivers(ctx context.Context, available bool) ([]*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers
		WHERE available = $1
	`
	var drivers []*driver.Driver
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &drivers, query, available)
	return drivers, err
}

func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM drivers WHERE id = $1`
	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify driver deletion: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("driver already deleted or not found")
	}
	return nil
}

func (r *DriverRepository) GetNearestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, current_location, ST_Distance(current_location, $1) AS dist
		FROM drivers
		WHERE available = true
		AND ST_DWithin(current_location, $1, $2)
		ORDER BY current_location <-> $1
		LIMIT 1
	`
	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, pickup, maxDistance)
	return &d, err
}
