package order

import (
	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	StoreID uuid.UUID `json:"store_id" binding:"required"`
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`

	Items []CreateOrderItem `json:"items" binding:"required,min=1"`

	PickupAddress string  `json:"pickup_address" binding:"required"`
	PickupLat     float64 `json:"pickup_lat" binding:"required"`
	PickupLng     float64 `json:"pickup_lng" binding:"required"`

	DeliveryAddress string  `json:"delivery_address" binding:"required"`
	DeliveryLat     float64 `json:"delivery_lat" binding:"required"`
	DeliveryLng     float64 `json:"delivery_lng" binding:"required"`
}

type CreateOrderItem struct {
	ProductID uuid.UUID  `json:"product_id" binding:"required"`
	VariantID *uuid.UUID `json:"variant_id,omitempty"`
	Quantity  int        `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderRequest struct {
	Column string      `json:"column" binding:"required"` // e.g. "status", "quantity"
	Value  interface{} `json:"value" binding:"required"`  // Accepts string, int, etc.
}

// ToOrder maps basic request info; snapshot fields filled in UseCase
func (r *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		StoreID:         r.StoreID,
		CustomerID:      r.CustomerID,
		PickupAddress:   r.PickupAddress,
		DeliveryAddress: r.DeliveryAddress,
		PickupPoint: postgis.PointS{
			X: r.PickupLng,
			Y: r.PickupLat,
			SRID: 4326,
		},
		DeliveryPoint: postgis.PointS{
			X: r.DeliveryLng,
			Y: r.DeliveryLat,
			SRID: 4326,
		},
		// ProductID, VariantID, Quantity, UnitPrice, Total, ProductName, VariantName, ImageURL
		// will be populated per item in UseCase
	}
}

// DropdownDataRequest represents the data used to populate order form dropdowns.
// swagger:model
type DropdownDataRequest struct {
	Customers   []Customer  `json:"customers"`
	Inventories []Inventory `json:"inventories"`
}

type Customer struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"full_name" json:"name"`
}

type Inventory struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	AdminID  uuid.UUID `db:"admin_id" json:"admin_id"`
	Category string    `db:"category" json:"category"`
}

// Point represents a simple GeoJSON-style point for Swagger only.
// swagger:model Point
type CreateOrderRequestDoc struct {
	MerchantID uuid.UUID `json:"merchant_id" binding:"required"`
	AdminID    uuid.UUID `json:"admin_id" binding:"required"`
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`

	ProductID uuid.UUID  `json:"product_id" binding:"required"`
	VariantID *uuid.UUID `json:"variant_id" binding:"required"`

	Quantity int `json:"quantity" binding:"required"`

	UnitPrice int64  `json:"unit_price" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
	Total     int64  `json:"total" binding:"required"`

	// snapshot at purchase time
	ProductName string `json:"product_name" binding:"required"`
	VariantName string `json:"variant_name" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required"`

	PickupAddress   string `json:"pickup_address" binding:"required"`
	PickupPoint     Point  `json:"pickup_point" binding:"required"`
	DeliveryAddress string `json:"delivery_address" binding:"required"`
	DeliveryPoint   Point  `json:"delivery_point" binding:"required"`
}
