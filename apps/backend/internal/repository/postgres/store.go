package postgres

import (
	"backend/internal/application"
	"backend/internal/domain/store"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type StoreRepository struct {
	exec sqlx.ExtContext
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{exec: db}
}

func (r *StoreRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *StoreRepository) Create(ctx context.Context, s *store.Store) error {
	query := `
    INSERT INTO stores (owner_id, name, name_normalized, location, logo_url)
		VALUES (:owner_id, :name, :name_normalized, :location, :logo_url)
		RETURNING id, created_at, updated_at
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, s)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return store.ErrStoreNameConflict
			case "23502":
				return store.ErrInvalidStoreInput
			}
		}
		return err
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return fmt.Errorf("scanning new store id: %w", err)
		}
	}

	return nil
}

func (r *StoreRepository) GetByID(ctx context.Context, storeID uuid.UUID) (*store.Store, error) {
	query := `
		SELECT * FROM stores 
		WHERE id =  $1
	`

	var s store.Store
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &s, query, storeID)
	return &s, err
}

func (r *StoreRepository) GetBasicByID(ctx context.Context, storeID uuid.UUID) (*store.StoreBasic, error) {
	query := `
		SELECT id, name, logo_url, rating FROM stores 
		WHERE id = $1
	`

	var s store.StoreBasic
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &s, query, storeID)
	return &s, err
}

func (r *StoreRepository) GetStoreSummary(ctx context.Context, storeID uuid.UUID) (*store.StoreSummary, error) {
	query := `
		SELECT id, name, logo_url, location, rating
		FROM stores
		WHERE id = $1
	`

	var summary store.StoreSummary
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &summary, query, storeID)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *StoreRepository) UpdateStoreDetails(ctx context.Context, storeID uuid.UUID, name string, logo string, location string) error {
	params := map[string]interface{}{
		"store_id": storeID,
		"name":     name,
		"logo_url": logo,
		"location": location,
	}

	query := `
		UPDATE stores
		SET name = :name,
			logo_url = :logo_url,
			location = :location,
			updated_at = NOW()
		WHERE id = :store_id
	`
	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23502": // not_null_violation
				return store.ErrInvalidStoreInput
			case "23505": // unique_violation
				return store.ErrStoreNameConflict
			case "23503": // foreign_key_violation
				return store.ErrInvalidStoreInput
			}
		}
		return fmt.Errorf("%w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return store.ErrStoreNotFound
	}

	return nil
}

func (r *StoreRepository) ListStores(ctx context.Context) ([]*store.Store, error) {
	query := `
		SELECT * FROM stores
	`

	var stores []*store.Store
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &stores, query)
	return stores, err
}

func (r *StoreRepository) ListStoresByOwner(ctx context.Context, ownerID uuid.UUID) ([]*store.MyStores, error) {
	query := `
		SELECT id, name FROM stores
		WHERE owner_id = $1
	`

	var stores []*store.MyStores
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &stores, query, ownerID)
	return stores, err
}

func (r *StoreRepository) ListStoresPaged(ctx context.Context, filter store.StoreFilter) ([]*store.StoreSummary, error) {
	query := `
		SELECT id, name, logo_url, location, rating
		FROM stores
		WHERE (:owner_id IS NULL OR owner_id = :owner_id)
		ORDER BY created_at DESC
		LIMIT :limit
		OFFSET :offset
	`

	params := map[string]interface{}{
		"owner_id": filter.OwnerID,
		"limit":    filter.Limit,
		"offset":   filter.Offset,
	}

	var stores []*store.StoreSummary
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &stores, query, params)
	return stores, err
}

func (r *StoreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM stores
		WHERE id = $1
	`

	res, err := r.execFromCtx(ctx).ExecContext(ctx, query, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				return store.ErrStoreHasReferences
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify store deletion: %w", err)
	}

	if rows == 0 {
		return store.ErrStoreNotFound
	}

	return nil
}

func (r *StoreRepository) Exists(ctx context.Context, storeID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM stores WHERE id = $1
		)
	`

	var exists bool
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &exists, query, storeID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *StoreRepository) IsOwnedBy(ctx context.Context, storeID uuid.UUID, ownerID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM stores
			WHERE id = $1 AND owner_id = $2
		)
	`

	var owned bool
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &owned, query, storeID, ownerID)
	if err != nil {
		return false, err
	}
	return owned, nil
}
