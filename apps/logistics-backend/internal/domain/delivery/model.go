package delivery

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryAssigned  DeliveryStatus = "assigned"
	DeliveryPickedUp  DeliveryStatus = "picked_up"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
)

type Delivery struct {
	ID          uuid.UUID      `db:"id" json:"id"`
	OrderID     uuid.UUID      `db:"order_id" json:"order_id"`
	DriverID    uuid.UUID      `db:"driver_id" json:"driver_id"`
	AssignedAt  time.Time      `db:"assigned_at" json:"assigned_at,omitzero"`
	PickedUpAt  time.Time      `db:"picked_up_at" json:"picked_up_at,omitzero"`
	DeliveredAt time.Time      `db:"delivered_at" json:"delivered_at,omitzero"`
	Status      DeliveryStatus `db:"status" json:"status"`
}
