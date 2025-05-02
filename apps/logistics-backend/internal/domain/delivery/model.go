package delivery

type DeliveryStatus string

const (
	DeliveryAssigned  DeliveryStatus = "assigned"
	DeliveryPickedUp  DeliveryStatus = "picked_up"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed"
)

type Delivery struct {
	ID          int64          `json:"id"`
	OrderID     int64          `json:"order_id"`
	DriverID    int64          `json:"driver_id"`
	AssignedAt  string         `json:"assigned_at"`
	PickedUpAt  string         `json:"picked_up_at"`
	DeliveredAt string         `json:"delivered_at"`
	Status      DeliveryStatus `json:"status"`
}
