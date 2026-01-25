package product

import "errors"

var (
	ErrProductInvalidStore       = errors.New("Invalid store reference.")
	ErrInvalidProduct            = errors.New("Invalid product reference.")
	ErrInvalidVariant            = errors.New("Invalid variant reference.")
	ErrInvalidOption             = errors.New("Invalid option reference.")
	ErrInvalidOptionValue        = errors.New("Invalid product or option reference.")
	ErrInvalidOptionValueInput   = errors.New("Invalid option value input.")
	ErrInvalidImage              = errors.New("Invalid image reference.")
	ErrInvalidReorder            = errors.New("Invalid product or images references.")
	ErrImageNotFound             = errors.New("Image not found.")
	ErrProductAlreadyExists      = errors.New("Product already exists.")
	ErrInvalidProductInput       = errors.New("Invalid product data.")
	ErrProductNotFound           = errors.New("Product not found.")
	ErrOptionNotFound            = errors.New("Option not found.")
	ErrOptionValueNotFound       = errors.New("Option value not found.")
	ErrMissingOptValues          = errors.New("No option values provided.")
	ErrInvalidOptionInput        = errors.New("Invalid option data.")
	ErrOptionValueInUse          = errors.New("Option value is in use by existing variants.")
	ErrOptionInUse               = errors.New("Option is in use by existing variants.")
	ErrInvalidVariantInput       = errors.New("Invalid variant data.")
	ErrInvalidVariantOptValInput = errors.New("Invalid variant option data.")
	ErrVariantAlreadyExists      = errors.New("Variant already exists.")
	ErrVariantNotFound           = errors.New("Variant not found.")
)
