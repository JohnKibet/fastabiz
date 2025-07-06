package postgres

import (
	"context"
	"database/sql"
	"logistics-backend/internal/domain/inventory"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRespository(db *sqlx.DB) inventory.Repository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(i *inventory.Inventory) error {
	query := `
		INSERT INTO inventories 
		(admin_id, name, category, stock, price, images, unit, packaging, description, location, slug)
		VALUES (:admin_id, :name, :category, :stock, :price, :images, :unit, :packaging, :description, :location, :slug)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&i.ID, i)
}

func (r *InventoryRepository) GetByID(InventoryID uuid.UUID) ([]*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price, images, unit, packaging, description, location, slug FROM inventories WHERE id = $1
	`
	var inventories []*inventory.Inventory
	err := r.db.Select(&inventories, query, InventoryID)
	return inventories, err
}

func (r *InventoryRepository) GetByName(InventoryName string) ([]*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price, images, unit, packaging, description, location, slug 
		FROM inventories 
		WHERE name = $1
	`
	var inventories []*inventory.Inventory
	err := r.db.Select(&inventories, query, InventoryName)
	if err != nil {
		return nil, err
	}
	if len(inventories) == 0 {
		return nil, sql.ErrNoRows
	}

	return inventories, nil
}

func (r *InventoryRepository) GetByCategory(ctx context.Context, category string) ([]inventory.Inventory, error) {
	query := `SELECT * FROM inventories WHERE category = $1`
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []inventory.Inventory
	for rows.Next() {
		var inv inventory.Inventory
		if err := rows.Scan(
			&inv.ID, &inv.AdminID, &inv.Name, &inv.Category,
			&inv.Stock, &inv.Price, &inv.Images, &inv.Unit,
			&inv.Packaging, &inv.Description, &inv.Location,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, nil
}

func (r *InventoryRepository) ListCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM inventories ORDER BY category`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *InventoryRepository) List(limit, offset int) ([]*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price, images, unit, packaging, description, location, slug, created_at, updated_at
		FROM inventories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	var inventories []*inventory.Inventory
	err := r.db.Select(&inventories, query, limit, offset)
	return inventories, err
}

func (r *InventoryRepository) GetBySlugs(adminSlug, productSlug string) (*inventory.Inventory, error) {
	query := `
		SELECT i.id, i.admin_id, i.name, i.category, i.stock, i.price, i.images, i.unit, i.packaging, i.description, i.location, i.slug
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE i.slug = $1 AND u.slug = $2
	`

	var inv inventory.Inventory
	err := r.db.Get(&inv, query, adminSlug, productSlug)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetStoreView(adminSlug string) (*inventory.StorePublicView, error) {

	// Getting admin info
	var store inventory.StorePublicView
	adminQuery := `
		SELECT full_name AS admin_name, category, location
		FROM users 
		WHERE slug = $1 
		AND role = 'admin'
	`
	if err := r.db.Get(&store, adminQuery, adminQuery); err != nil {
		return nil, err
	}

	// Getting products for this admin
	productQuery := `
		SELECT name, price, unit, packaging, stock AS in_stock,
			(split_part(images, ',', 1)) AS image
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE u.slug = $1
	`
	err := r.db.Select(&store.Products, productQuery, adminQuery)
	if err != nil {
		return nil, err
	}

	return &store, nil
}
