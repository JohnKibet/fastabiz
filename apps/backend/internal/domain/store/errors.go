package store

import "errors"

var (
	ErrCreateStore            = errors.New("Create store failed.")
	ErrNotOwner               = errors.New("User does not own store.")
	ErrStoreNotFound          = errors.New("Store not found.")
	ErrStoreHasReferences     = errors.New("Store has dependent records.")
	ErrInvalidStoreInput      = errors.New("Invalid store data.")
	ErrStoreNameConflict      = errors.New("Store name already exists.")
)
