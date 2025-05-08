package delivery

import (
	"github.com/google/uuid"
)

type CreateDeliveryRequest struct {
	OrderID  uuid.UUID      `json:"order_id"`
	DriverID uuid.UUID      `json:"driver_id"`
	Status   DeliveryStatus `json:"status"`
}

func (r *CreateDeliveryRequest) ToDelivery() *Delivery {
	return &Delivery{
		OrderID:  r.OrderID,
		DriverID: r.DriverID,
		Status:   DeliveryAssigned,
	}
}
