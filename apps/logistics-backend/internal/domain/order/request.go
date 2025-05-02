package order

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID       uuid.UUID `json:"customer_id" binding:"required"`
	PickupLocation   string    `json:"pickup_location" binding:"required"`
	DeliveryLocation string    `json:"delivery_location" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}

func (r *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		CustomerID:       r.CustomerID,
		PickupLocation:   r.PickupLocation,
		DeliveryLocation: r.DeliveryLocation,
		OrderStatus:      Pending,
	}
}

func (r *UpdateOrderStatusRequest) ToUpdateOrder() *Order {
	return &Order{
		OrderStatus: r.Status,
	}
}
