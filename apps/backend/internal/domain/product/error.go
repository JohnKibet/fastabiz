package product

import "errors"

var (
	ErrProductInvalidStore  = errors.New("invalid store reference")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductInvalidInput  = errors.New("invalid product data")
	ErrProductNotFound      = errors.New("product not found")
	ErrOptionNotFound       = errors.New("option not found")
	ErrOptionValueNotFound  = errors.New("option value not found")
	ErrOptionValueInUse     = errors.New("option value is in use by existing variants")
	ErrOptionInUse          = errors.New("option is in use by existing variants")
)
