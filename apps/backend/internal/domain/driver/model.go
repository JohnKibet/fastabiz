package driver

import (
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type Driver struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	FullName        string         `db:"full_name" json:"full_name"`
	Email           string         `db:"email" json:"email"`
	VehicleInfo     string         `db:"vehicle_info" json:"vehicle_info"`
	CurrentLocation postgis.PointS `db:"current_location" json:"current_location"`
	Available       bool           `db:"available" json:"available"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
}

// Point represents a simple GeoJSON-style point for Swagger only.
// swagger:model Point
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DriverDoc is used only for Swagger documentation.
// swagger:model DriverDoc
type DriverDoc struct {
	ID              uuid.UUID `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	VehicleInfo     string    `json:"vehicle_info"`
	CurrentLocation Point     `json:"current_location"`
	Available       bool      `json:"available"`
	CreatedAt       time.Time `json:"created_at"`
}
