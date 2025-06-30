package inventory

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID          uuid.UUID `db:"id" json:"id"`
	AdminID     uuid.UUID `db:"admin_id" json:"admin_id"` // Foreign key
	Name        string    `db:"name" json:"name"`         // e.g. “Fresh Milk”
	Category    string    `db:"category" json:"category"` // e.g. “Dairy”
	Stock       int       `db:"stock" json:"stock"`
	Price       float64   `db:"price" json:"price"`
	Images      string    `db:"images" json:"images"`       // could be JSON array or URLs
	Unit        string    `db:"unit" json:"unit"`           // "per litre", "per bucket"
	Packaging   string    `db:"packaging" json:"packaging"` // “Bucket/Single”
	Description string    `db:"description" json:"description"`
	Location    string    `db:"location" json:"location"` // optional
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type StorePublicView struct {
	AdminName string             `json:"admin_name"`
	Category  string             `json:"category"`
	Location  string             `json:"location"`
	Products  []InventorySummary `json:"products"`
}

type InventorySummary struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
	Unit      string  `json:"unit"`
	Packaging string  `json:"packaging"`
	InStock   int     `json:"in_stock"`
}
