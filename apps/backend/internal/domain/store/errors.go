package store

import "errors"

var ErrStoreNameAlreadyExists = errors.New("store name already exists for this owner")
