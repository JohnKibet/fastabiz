package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `db:"id" json:"id"`
	MerchantID  uuid.UUID `db:"merchant_id" json:"merchant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Category    string    `db:"category" json:"category"`

	// Stored as JSON/JSONB array in DB
	Images []string `db:"images" json:"images"`

	HasVariants bool `db:"has_variants" json:"has_variants"`

	// Used only when HasVariants == false
	Price float64 `db:"price" json:"price,omitempty"`
	Stock int     `db:"stock" json:"stock,omitempty"`

	// Stored as JSON/JSONB
	Options  []Option  `db:"options" json:"options,omitempty"`
	Variants []Variant `db:"variants" json:"variants,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Option struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Name   string    `db:"name" json:"name"`
	Values []string  `db:"values" json:"values"`
}

type Variant struct {
	ID        uuid.UUID               `db:"id" json:"id"`
	ProductID uuid.UUID               `db:"product_id" json:"product_id"`
	SKU       string                  `db:"sku" json:"sku"`
	Price     float64                 `db:"price" json:"price"`
	Stock     int                     `db:"stock" json:"stock"`
	ImageURL  string                  `db:"image_url" json:"image_url"`
	Options   map[uuid.UUID]uuid.UUID `db:"options" json:"options"`
}
