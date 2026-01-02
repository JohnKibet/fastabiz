package store

import "github.com/google/uuid"

type CreateStoreRequest struct {
	OwnerID  uuid.UUID `db:"admin_id" binding:"required"`
	Name     string    `db:"name" binding:"required"` // example:"Kevin's Electronics"
	LogoURL  string    `json:"logo_url"`              // example:"https://cdn.fastabiz.com/logos/kevins.png"
	Location string    `json:"location" binding:"required"`
}

type UpdateStoreRequest struct {
	Name     string
	Location string
	Logo     string
}

func (r *CreateStoreRequest) ToStore() *Store {
	return &Store{
		OwnerID:  r.OwnerID,
		Name:     r.Name,
		LogoURL:  r.LogoURL,
		Location: r.Location,
	}
}
