package store

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create persists a new store aggregate, including its core identifying and ownership attributes.

	// Child or derived data (analytics, summaries, statistics)
	// are expected to be populated or updated separately.
	Create(ctx context.Context, s *Store) error

	// UpdateStoreDetails updates the mutable core attributes of a store,
	// such as its display name, branding, and physical or business location.
	//
	// This does not affect ownership, permissions, or derived analytics.
	UpdateStoreDetails(ctx context.Context, storeID uuid.UUID, name, logo, location string) error

	// Delete permanently removes a store and all associated child records,
	// such as products, configurations, and dependent metadata.
	Delete(ctx context.Context, storeID uuid.UUID) error

	// GetByID retrieves a fully-hydrated store aggregate by its ID.
	//
	// This is intended for internal domain use where complete store
	// state is required.
	GetByID(ctx context.Context, storeID uuid.UUID) (*Store, error)

	// GetBasicByID retrieves a lightweight projection of a store,
	// containing only essential fields required for references,
	// headers, or authorization checks.
	GetBasicByID(ctx context.Context, storeID uuid.UUID) (*StoreBasic, error)

	// ListStores returns all stores accessible to the caller.
	//
	// Each store should be returned as a full aggregate and may span
	// multiple owners, depending on caller permissions.
	ListStores(ctx context.Context) ([]*Store, error)

	// ListStoresByOwner returns all stores owned by the specified owner.
	//
	// This is typically used for dashboards and management views
	// scoped to a single account.
	ListStoresByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Store, error)

	// ListStoresPaged returns a paginated and filtered list of stores
	// using the provided filter criteria.
	//
	// Results are returned as summarized projections suitable for
	// tables, search results, and administrative listings.
	ListStoresPaged(ctx context.Context, filter StoreFilter) ([]*StoreSummary, error)

	// GetStoreSummary retrieves an aggregated, read-optimized view
	// of a store, including derived metrics such as product counts,
	// order statistics, or revenue indicators.
	//
	// This method does not return the full store aggregate.
	GetStoreSummary(ctx context.Context, storeID uuid.UUID) (*StoreSummary, error)

	// Exists checks whether a store with the given ID exists.
	//
	// This is a lightweight guard method and should not load
	// the full store aggregate.
	Exists(ctx context.Context, storeID uuid.UUID) (bool, error)

	// IsOwnedBy verifies whether the specified store is owned
	// by the given owner.
	//
	// This is primarily used for authorization and access control
	// checks prior to executing write operations.
	IsOwnedBy(ctx context.Context, storeID uuid.UUID, ownerID uuid.UUID) (bool, error)
}
