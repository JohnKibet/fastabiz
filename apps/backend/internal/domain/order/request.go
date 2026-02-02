package order

import (
	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	StoreID    uuid.UUID `json:"store_id" binding:"required"`
	AdminID    uuid.UUID `json:"admin_id"`
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

	PickupAddress   string         `json:"pickup_address" binding:"required"`
	PickupPoint     postgis.PointS `json:"pickup_point" binding:"required"`
	DeliveryAddress string         `json:"delivery_address" binding:"required"`
	DeliveryPoint   postgis.PointS `json:"delivery_point" binding:"required"`
}

type UpdateOrderRequest struct {
	Column string      `json:"column" binding:"required"` // e.g. "status", "quantity"
	Value  interface{} `json:"value" binding:"required"`  // Accepts string, int, etc.
}

func (r *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		StoreID:    r.StoreID,
		AdminID:    r.AdminID,
		CustomerID: r.CustomerID,

		ProductID: r.ProductID,
		VariantID: r.VariantID,

		Quantity:  r.Quantity,
		UnitPrice: r.UnitPrice,
		Currency:  r.Currency,
		Total:     r.Total,

		ProductName: r.ProductName,
		VariantName: r.VariantName,
		ImageURL:    r.ImageURL,

		PickupAddress:   r.PickupAddress,
		DeliveryAddress: r.DeliveryAddress,
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
