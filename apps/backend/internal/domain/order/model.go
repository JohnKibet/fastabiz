package order

import (
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Assigned  OrderStatus = "assigned"
	InTransit OrderStatus = "in_transit"
	Delivered OrderStatus = "delivered"
	Cancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         uuid.UUID `db:"id" json:"id"`
	MerchantID uuid.UUID `db:"merchant_id" json:"merchant_id"` // owner
	AdminID    uuid.UUID `db:"admin_id" json:"admin_id"`       // manager
	CustomerID uuid.UUID `db:"user_id" json:"customer_id"`

	ProductID uuid.UUID  `db:"product_id" json:"product_id"`
	VariantID *uuid.UUID `db:"variant_id" json:"variant_id"` // NULLABLE

	Quantity int `db:"quantity" json:"quantity"`

	// Price snapshot
	UnitPrice int64  `db:"unit_price" json:"unit_price"`
	Currency  string `db:"currency" json:"currency"`
	Total     int64  `db:"total" json:"total"` // quantity * unit_price

	// Optional — snapshot of name & image at purchase time
	ProductName string `db:"product_name" json:"product_name"`
	VariantName string `db:"variant_name" json:"variant_name"`
	ImageURL    string `db:"image_url" json:"image_url"`

	// Pickup/delivery
	PickupAddress   string         `db:"pickup_address" json:"pickup_address"`
	PickupPoint     postgis.PointS `db:"pickup_point" json:"pickup_point"`
	DeliveryAddress string         `db:"delivery_address" json:"delivery_address"`
	DeliveryPoint   postgis.PointS `db:"delivery_point" json:"delivery_point"`

	Status    OrderStatus `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

// Point represents a simple GeoJSON-style point for Swagger only.
// swagger:model Point
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type OrderDoc struct {
	ID         uuid.UUID `db:"id" json:"id"`
	MerchantID uuid.UUID `db:"merchant_id" json:"merchant_id"`
	AdminID    uuid.UUID `db:"admin_id" json:"admin_id"`
	CustomerID uuid.UUID `db:"user_id" json:"customer_id"`

	ProductID uuid.UUID  `db:"product_id" json:"product_id"`
	VariantID *uuid.UUID `db:"variant_id" json:"variant_id"`

	Quantity int `db:"quantity" json:"quantity"`

	UnitPrice int64  `db:"unit_price" json:"unit_price"`
	Currency  string `db:"currency" json:"currency"`
	Total     int64  `db:"total" json:"total"`

	ProductName string `db:"product_name" json:"product_name"`
	VariantName string `db:"variant_name" json:"variant_name"`
	ImageURL    string `db:"image_url" json:"image_url"`

	PickupAddress   string `db:"pickup_address" json:"pickup_address"`
	PickupPoint     Point  `db:"pickup_point" json:"pickup_point"`
	DeliveryAddress string `db:"delivery_address" json:"delivery_address"`
	DeliveryPoint   Point  `db:"delivery_point" json:"delivery_point"`

	Status    OrderStatus `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}
