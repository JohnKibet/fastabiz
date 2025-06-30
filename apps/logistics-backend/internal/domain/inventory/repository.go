package inventory

import (
	"github.com/google/uuid"
)

type Repository interface {
	Create(inventory *Inventory) error
	GetByID(id uuid.UUID) ([]*Inventory, error)
	GetByName(name string) ([]*Inventory, error)
	List(limit, offset int) ([]*Inventory, error)

	GetBySlugs(adminSlugs, productSlug string) (*Inventory, error)
	GetStoreView(adminSlug string) (*StorePublicView, error)
}
