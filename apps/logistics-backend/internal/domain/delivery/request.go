package delivery

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeliveryRequest struct {
	OrderID  uuid.UUID      `json:"order_id"`
	DriverID uuid.UUID      `json:"driver_id"`
	Status   DeliveryStatus `json:"status"`
}

type UpdateDeliveryRequest struct {
	Column string      `json:"column" binding:"required"`
	Value  interface{} `json:"value" binding:"required"`
}

func (r *CreateDeliveryRequest) ToDelivery() *Delivery {
	return &Delivery{
		OrderID:    r.OrderID,
		DriverID:   r.DriverID,
		Status:     DeliveryAssigned,
		AssignedAt: time.Now(),
	}
}
