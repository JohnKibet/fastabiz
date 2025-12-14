package order

import (
	"backend/internal/domain/driver"
	"backend/internal/domain/inventory"
	"backend/internal/domain/notification"
	"backend/internal/domain/product"
	"backend/internal/domain/store"
	"context"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

// cross-domain DI using necessary interface

// Access the inventory domain usecase methods.
type InventoryReader interface {
	GetInventoryByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error)
	UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error
	GetAllInventories(ctx context.Context) ([]Inventory, error)
}

// Access the user domain usecase method for getting users of role customers.
type CustomerReader interface {
	GetAllCustomers(ctx context.Context) ([]Customer, error)
}

type DriverReader interface {
	GetNearestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*driver.Driver, error)
}

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}

type StoreReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*store.Store, error)
}

// new interfaces to replace store and inventory readers
type ProductOrVariantReader interface {
	GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error)
	GetVariantByID(ctx context.Context, id uuid.UUID) (*product.Variant, error)
	UpdateProductStock(ctx context.Context, productID uuid.UUID, newStock int) error
	UpdateVariantStock(ctx context.Context, variantID uuid.UUID, newStock int) error
}
