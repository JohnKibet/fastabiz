package store

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OwnerID        uuid.UUID `db:"owner_id" json:"admin_id"` // FK to users
	Name           string    `db:"name" json:"name"`
	NameNormalized string    `db:"name_normalized" json:"-"`
	LogoURL        string    `db:"logo_url" json:"logo_url"`
	Location       string    `db:"location" json:"location"` // optional
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// Store analytics (derived data)
// This is used for:
// Store profile page
// Merchant dashboard
// Admin views
type StoreSummary struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	LogoURL  string    `db:"logo_url"`
	Rating   float64   `db:"rating"`
	Location string    `db:"location"`
}

// GetBasicByID(ctx context.Context, storeID uuid.UUID) (*StoreBasic, error)
// This avoids:
// Loading full store rows
// Repeating JOIN logic across repositories
type StoreBasic struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Logo   string    `json:"logo_url"`
	Rating float64   `json:"rating"`
}

// Marketplace listing (paged + ordered)
// Used for:
// Marketplace home
// “Browse stores”
// Sorting by rating or popularity
type StoreFilter struct {
	OwnerID *uuid.UUID `db:"owner_id" json:"owner_id"`
	Limit   int        `db:"limit" json:"limit"`
	Offset  int        `db:"offset" json:"offset"`
}

// Existence / ownership checks (important for authorization)
// Used for:
// “Can this user edit this store?”
// “Can this merchant create products under this store?”
