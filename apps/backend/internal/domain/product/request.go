package product

import "github.com/google/uuid"

type CreateVariantRequest struct {
	ProductID uuid.UUID               `json:"product_id" binding:"required"`
	SKU       string                  `json:"sku" binding:"required"`
	Price     float64                 `json:"price" binding:"required"`
	Stock     int                     `json:"stock" binding:"required"`
	ImageURL  *string                 `json:"image_url" binding:"required"`
	Options   map[uuid.UUID]uuid.UUID `json:"options" binding:"required"` // optionID -> optionValueID
}

type CreateProductRequest struct {
	StoreID     uuid.UUID `json:"store_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Category    string    `json:"category" binding:"required"`
}

type UpdateProductDetailsRequest struct {
	ProductID   uuid.UUID `json:"product_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Category    string    `json:"category" binding:"required"`
}

type AddImageRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	URL       string    `json:"url" binding:"required"`
	IsPrimary bool      `json:"is_primary"`
}

type ReorderImagesRequest struct {
	ProductID uuid.UUID   `json:"product_id" binding:"required"`
	ImageIDs  []uuid.UUID `json:"image_ids" binding:"required"`
}

type AddOptionNameRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}

type AddOptionValueRequest struct {
	OptionID uuid.UUID `json:"option_id" binding:"required"`
	Value    string    `json:"value" binding:"required"`
}

type UpdateVariantStockRequest struct {
	VariantID uuid.UUID `json:"variant_id" binding:"required"`
	Stock     int       `json:"stock" binding:"required"`
}

type UpdateVariantPriceRequest struct {
	VariantID uuid.UUID `json:"variant_id" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}

func (r CreateProductRequest) ToProduct() *Product {
	return &Product{
		StoreID:     r.StoreID,
		Name:        r.Name,
		Description: r.Description,
		Category:    r.Category,
	}
}

func (r CreateVariantRequest) ToVariant() *Variant {
	return &Variant{
		ProductID: r.ProductID,
		SKU:       r.SKU,
		Price:     r.Price,
		Stock:     r.Stock,
		ImageURL:  *r.ImageURL,
		Options:   r.Options,
	}
}
