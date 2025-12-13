package product

import "github.com/google/uuid"

type VariantInput struct {
	ID        string                  `db:"id" json:"id"`
	ProductID uuid.UUID               `json:"product_id" binding:"required"`
	SKU       string                  `json:"sku" binding:"required"`
	Price     float64                 `json:"price" binding:"required"`
	Stock     int                     `json:"stock" binding:"required"`
	ImageURL  *string                 `json:"image_url" binding:"required"`
	Options   map[uuid.UUID]uuid.UUID `json:"options" binding:"required"` // optionID -> optionValueID
}
