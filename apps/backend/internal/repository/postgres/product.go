package postgres

import (
	"backend/internal/application"
	"backend/internal/domain/product"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductRepository struct {
	exec sqlx.ExtContext
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{exec: db}
}

func (r *ProductRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *ProductRepository) AddImage(ctx context.Context, productID uuid.UUID, url string, isPrimary bool) error {
	params := map[string]interface{}{
		"product_id": productID,
		"url":        url,
		"is_primary": isPrimary,
	}

	query := `
		INSERT INTO product_images (product_id, url, is_primary)
        VALUES (:product_id, :url, :is_primary)
        RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("insert image: %w", err)
	}
	defer rows.Close()

	var imageID uuid.UUID
	if rows.Next() {
		if err := rows.Scan(&imageID); err != nil {
			return fmt.Errorf("scanning new image id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after insert")
	}

	return nil
}

func (r *ProductRepository) RemoveImage(ctx context.Context, imageID uuid.UUID) error {
	query := `
		DELETE FROM product_images
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify image deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("image already deleted or invalid")
	}

	return nil
}

func (r *ProductRepository) ReorderImages(ctx context.Context, productID uuid.UUID, imageIDs []uuid.UUID) error {
	if len(imageIDs) == 0 {
		return nil
	}

	params := map[string]interface{}{
		"product_id": productID,
		"image_ids":  pq.Array(imageIDs),
	}

	query := `
        WITH ordered_images AS (
            SELECT
                unnest(:image_ids::uuid[]) AS image_id,
                ordinality - 1 AS position
            FROM unnest(:image_ids::uuid[]) WITH ORDINALITY
        )
        UPDATE product_images pi
        SET position = oi.position
        FROM ordered_images oi
        WHERE pi.id = oi.image_id
        AND pi.product_id = :product_id
    `

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("reorder product images: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows != int64(len(imageIDs)) {
		return fmt.Errorf(
			"image reorder mismatch: expected %d updates, got %d",
			len(imageIDs), rows,
		)
	}

	return nil
}

func (r *ProductRepository) AddOption(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error) {
	params := map[string]interface{}{
		"product_id": productID,
		"name":       name,
	}

	query := `
		INSERT INTO product_options (product_id, name)
        VALUES (:product_id, :name)
        RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert option name: %w", err)
	}
	defer rows.Close()

	var optionID uuid.UUID
	if rows.Next() {
		if err := rows.Scan(&optionID); err != nil {
			return uuid.Nil, fmt.Errorf("scanning new option id: %w", err)
		}
	} else {
		return uuid.Nil, fmt.Errorf("no id returned after insert")
	}

	return optionID, nil
}

func (r *ProductRepository) RemoveOption(ctx context.Context, optionID uuid.UUID) error {
	query := `
		DELETE FROM product_options
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, optionID)
	if err != nil {
		return fmt.Errorf("failed to delete option name: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify option name deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("option name already deleted or invalid")
	}

	return nil
}

func (r *ProductRepository) AddOptionValue(ctx context.Context, optionID uuid.UUID, value string) error {
	params := map[string]interface{}{
		"option_id": optionID,
		"value":     value,
	}

	query := `
		INSERT INTO product_option_values (option_id, value)
        VALUES (:option_id, :value)
        RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("insert option value: %w", err)
	}
	defer rows.Close()

	var optionValueID uuid.UUID
	if rows.Next() {
		if err := rows.Scan(&optionValueID); err != nil {
			return fmt.Errorf("scanning new option value id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after insert")
	}

	return nil
}

func (r *ProductRepository) RemoveOptionValue(ctx context.Context, optionValueID uuid.UUID) error {
	query := `
		DELETE FROM product_option_values
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, optionValueID)
	if err != nil {
		return fmt.Errorf("failed to delete option value: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify option value deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("option value already deleted or invalid")
	}

	return nil
}

func (r *ProductRepository) CreateVariant(ctx context.Context, variant *product.Variant) error {
	query := `
		INSERT INTO variants (
			product_id, sku, price, stock, image_url, options
		) VALUES (
			:product_id, :sku, :price, :stock, :image_url, :options
		)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, variant)
	if err != nil {
		return fmt.Errorf("insert variant: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&variant.ID); err != nil {
			return fmt.Errorf("scanning new variant id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *ProductRepository) UpdateVariantStock(ctx context.Context, variantID uuid.UUID, stock int) error {
	params := map[string]interface{}{
		"variant_id": variantID,
		"stock":      stock,
	}

	query := `
		UPDATE variants
        SET stock = :stock,
            updated_at = NOW()
        WHERE id = :variant_id
	`

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("update variant stock: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no variant found with id %s", variantID)
	}

	return nil
}

func (r *ProductRepository) UpdateVariantPrice(ctx context.Context, variantID uuid.UUID, price float64) error {
	params := map[string]interface{}{
		"variant_id": variantID,
		"price":      price,
	}

	query := `
		UPDATE variants
        SET price = :price,
            updated_at = NOW()
        WHERE id = :variant_id
	`

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("update variant price: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no variant found with id %s", variantID)
	}

	return nil
}

func (r *ProductRepository) RemoveVariant(ctx context.Context, variantID uuid.UUID) error {
	query := `
		DELETE FROM variants 
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, variantID)
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify variant deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("variant already deleted or invalid")
	}

	return nil
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (
		merchant_id, name, description, category
	) VALUES (
		:merchant_id, :name, :description, :category
	)
	RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, p)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.ID); err != nil {
			return fmt.Errorf("scanning new product id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after insert")
	}

	// 2. Insert images (optional)
	for _, img := range p.Images {
		if err := r.AddImage(ctx, p.ID, img, false); err != nil {
			return err
		}
	}

	// 3. Insert options and option values (optional)
	for _, opt := range p.Options {
		// insert option and get its ID
		optionID, err := r.AddOption(ctx, p.ID, opt.Name)
		if err != nil {
			return err
		}

		// insert each value
		for _, val := range opt.Values {
			if err := r.AddOptionValue(ctx, optionID, val); err != nil {
				return err
			}
		}
	}

	// 4. Insert variants if HasVariants == true
	if p.HasVariants {
		for _, v := range p.Variants {
			// mapping - mv top level later
			vi := &product.Variant{
				ID:        p.ID,
				ProductID: p.ID,
				SKU:       v.SKU,
				Price:     v.Price,
				Stock:     v.Stock,
				Options:   v.Options,
			}

			if v.ImageURL != "" {
				vi.ImageURL = v.ImageURL
			}

			// insert into variants and variant_option_values
			if err := r.CreateVariant(ctx, vi); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := `
		SELECT id, merchant_id, name, description, category
		FROM products 
		WHERE id = $1
	`

	var p product.Product
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &p, query, id); err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &p, nil
}

func (r *ProductRepository) UpdateDetails(ctx context.Context, productID uuid.UUID, name, description, category string) error {
	params := map[string]interface{}{
		"product_id":  productID,
		"name":        name,
		"description": description,
		"category":    category,
	}

	query := `
		UPDATE products
		SET name = :name,
			description = :description,
			category = :category,
			updated_at = NOW()
		WHERE id = :product_id
	`
	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		return fmt.Errorf("update product details: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no product found with id %s", productID)
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context) ([]product.Product, error) {
	query := `
		SELECT id, merchant_id, name, description, category
		FROM products
	`

	var products []product.Product
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &products, query)
	return products, err
}

func (r *ProductRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	query := `
		DELETE FROM products 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify product deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product already deleted or invalid")
	}

	return nil
}
