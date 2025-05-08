package driver

import "github.com/google/uuid"

type Driver struct {
	ID              uuid.UUID `db:"id" json:"id"`
	FullName        string    `db:"full_name" json:"full_name"`
	Email           string    `db:"email" json:"email"`
	VehicleInfo     string    `db:"vehicle_info" json:"vehicle_info"`
	CurrentLocation string    `db:"current_location" json:"current_location"`
	Available       bool      `db:"available" json:"available"`
	CreatedAt       string    `db:"created_at" json:"created_at"`
}
