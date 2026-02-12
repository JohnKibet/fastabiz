package product

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {

	// Create persists a new product aggregate, including its core attributes.
	// Child entities such as images, options, and variants may be created
	// separately using their respective methods.
	Create(ctx context.Context, product *Product) error

	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)

	// GetProductByID retrieves a fully-hydrated product aggregate by its ID,
	// including images, options, option values, and variants.
	GetFullProductByID(ctx context.Context, id uuid.UUID) (*Product, error)

	UpdateProductStock(ctx context.Context, productID uuid.UUID, stock int) error

	// List returns all products from a specified store accessible to the caller.
	// Each product should be returned as a fully-hydrated aggregate.
	ListProductsByStore(ctx context.Context, storeID uuid.UUID) ([]ProductListItem, error)

	// UpdateDetails updates the core mutable fields of a product.
	// This does not affect images, options, variants, or inventory.
	UpdateDetails(
		ctx context.Context,
		productID uuid.UUID,
		name string,
		description string,
		category string,
	) error

	// CreateVariant creates a new purchasable variant for a product
	// and associates it with the specified option values.
	CreateVariant(ctx context.Context, variant *Variant) error

	// AddVariantOptionValue associates an option-value pair with a specific variant.
	AddVariantOptionValue(ctx context.Context, variantID uuid.UUID, valueID uuid.UUID) error

	// GetVariantByID retrieves a specific variant by its ID,
	// including its associated option-value mappings.
	GetVariantByID(ctx context.Context, variantID uuid.UUID) (*Variant, error)

	// UpdateVariantStock updates the available inventory quantity
	// for a specific variant.
	UpdateVariantStock(ctx context.Context, variantID uuid.UUID, stock int) error

	// UpdateVariantPrice updates the selling price for a specific variant.
	UpdateVariantPrice(ctx context.Context, variantID uuid.UUID, price float64) error

	// DeleteVariant permanently removes a variant and all associated
	// option-value mappings.
	RemoveVariant(ctx context.Context, variantID uuid.UUID) error

	// AddImage attaches a new image to a product.
	// If isPrimary is true, existing primary images should be unset.
	AddImage(ctx context.Context, productID uuid.UUID, url string, isPrimary bool) error

	// RemoveImage deletes an image from a product by image ID.
	RemoveImage(ctx context.Context, imageID uuid.UUID) error

	// ReorderImages updates the display order of product images
	// using the provided ordered list of image IDs.
	ReorderImages(
		ctx context.Context,
		productID uuid.UUID,
		imageIDs []uuid.UUID,
	) error

	// AddOption creates a new configurable option for a product
	// (e.g., Size, Weight, Color).
	AddOption(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error)

	// RemoveOption deletes a product option and all of its associated values.
	RemoveOption(ctx context.Context, optionID uuid.UUID) error

	// AddOptionValue adds a selectable value to an existing product option
	// (e.g., Small, Medium, Large).
	AddOptionValue(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, value string) error

	// RemoveOptionValue deletes a specific option value and removes
	// any variant associations that depend on it.
	RemoveOptionValue(ctx context.Context, optionValueID uuid.UUID) error

	// Delete permanently removes a product and all associated child records,
	// including images, options, variants, and inventory.
	Delete(ctx context.Context, productID uuid.UUID) error

	// GetOptionIDByName retrieves the ID of a product option by its name.
	GetOptionIDByName(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error)

	// GetOptionValueID retrieves the ID of a product option value by its value.
	GetOptionValueID(ctx context.Context, optionID uuid.UUID, value string) (uuid.UUID, error)

	// ListVariantsByProductID retrieves all variants associated with a specific product.
	ListVariantsByProductID(ctx context.Context, productID uuid.UUID) ([]Variant, error)

	// ListOptionsByProductID retrieves all options and their values for a specific product.
	ListOptionsByProductID(ctx context.Context, productID uuid.UUID) ([]Option, error)

	// ListOptionValuesByOptionID retrieves all values for a specific product option.
	ListOptionValuesByOptionID(ctx context.Context, optionID uuid.UUID) ([]OptionValue, error)

	// IsOptionUsed checks whether a product option (e.g. "Size", "Color") is currently
	// used by any variant.
	//
	// NOTE:
	// Variants do not reference product options directly.
	// The relationship is indirect:
	//
	//   product_options
	//        ↓
	//   product_option_values
	//        ↓
	//   variant_option_values
	//
	// Because variants only store references to option *values*, we must JOIN
	// `variant_option_values` with `product_option_values` to determine whether
	// any variant uses *any value* belonging to the given option.
	//
	// This check is required to prevent deleting an option that is still in use
	// by one or more variants, which would otherwise leave orphaned variant data
	// and violate domain integrity.
	IsOptionUsed(ctx context.Context, optionID uuid.UUID) (bool, error)

	// IsOptionValueUsed checks whether a product option value (e.g. "Small", "Red") is currently
	// used by any variant.
	//
	// This check is required to prevent deleting an option value that is still in use
	// by one or more variants, which would otherwise leave orphaned variant data
	// and violate domain integrity.
	IsOptionValueUsed(ctx context.Context, optionValueID uuid.UUID) (bool, error)

	// Universal fetch for all products regardless of store or merchant
	// ListAllProducts()

	// Adds price to product without variants
	UpdateProductInvPrice(ctx context.Context, productID uuid.UUID, price float64) error

	// Adds stock to product without variants
	UpdateProductInvStock(ctx context.Context, productID uuid.UUID, stock int) error
}
