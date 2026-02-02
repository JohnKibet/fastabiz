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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidProduct
			}
		}
		return err
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidImage
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify image deletion: %w", err)
	}

	if rows == 0 {
		return product.ErrImageNotFound
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidReorder
			}
		}
		return err
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return uuid.Nil, product.ErrInvalidProduct
			case "23502":
				return uuid.Nil, product.ErrInvalidOptionInput
			}
		}
		return uuid.Nil, err
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

func (r *ProductRepository) GetOptionIDByName(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error) {
	query := `
		SELECT id FROM product_options
		WHERE product_id = $1 AND name = $2
	`

	var optionID uuid.UUID
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &optionID, query, productID, name); err != nil {
		return uuid.Nil, fmt.Errorf("get option id by name: %w", err)
	}

	return optionID, nil
}

func (r *ProductRepository) ListOptionsByProductID(ctx context.Context, productID uuid.UUID) ([]product.Option, error) {
	query := `
		SELECT id, name
		FROM product_options
		WHERE product_id = $1
		ORDER BY position
	`

	var options []product.Option
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &options, query, productID)
	return options, err
}

func (r *ProductRepository) RemoveOption(ctx context.Context, optionID uuid.UUID) error {
	query := `
		DELETE FROM product_options
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, optionID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidOption
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify option name deletion: %w", err)
	}

	if rows == 0 {
		return product.ErrOptionNotFound
	}

	return nil
}

func (r *ProductRepository) AddOptionValue(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, value string) error {
	params := map[string]interface{}{
		"product_option_id": optionID,
		"product_id":        productID,
		"value":             value,
	}

	query := `
		INSERT INTO product_option_values (product_option_id, product_id, value)
    VALUES (:product_option_id, :product_id, :value)
    RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidOptionValue
			case "23502":
				return product.ErrInvalidOptionValueInput
			}
		}
		return err
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

func (r *ProductRepository) GetOptionValueID(ctx context.Context, optionID uuid.UUID, value string) (uuid.UUID, error) {
	query := `
		SELECT id FROM product_option_values
		WHERE product_option_id = $1 AND value = $2
	`

	var optionValueID uuid.UUID
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &optionValueID, query, optionID, value); err != nil {
		return uuid.Nil, fmt.Errorf("get option value id: %w", err)
	}

	return optionValueID, nil
}

func (r *ProductRepository) ListOptionValuesByOptionID(ctx context.Context, optionID uuid.UUID) ([]product.OptionValue, error) {

	query := `
		SELECT id, value
		FROM product_option_values
		WHERE product_option_id = $1
		ORDER BY position
	`

	var values []product.OptionValue
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &values, query, optionID)
	return values, err
}

func (r *ProductRepository) RemoveOptionValue(ctx context.Context, optionValueID uuid.UUID) error {
	query := `
		DELETE FROM product_option_values
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, optionValueID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidOptionValueInput
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify option value deletion: %w", err)
	}

	if rows == 0 {
		return product.ErrOptionValueNotFound
	}

	return nil
}

func (r *ProductRepository) CreateVariant(ctx context.Context, variant *product.Variant) error {
	query := `
		INSERT INTO variants (
			product_id, sku, price, stock, image_url
		) VALUES (
			:product_id, :sku, :price, :stock, :image_url
		)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, variant)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return product.ErrVariantAlreadyExists
			case "23502": // not_null_violation
				return product.ErrInvalidVariantInput
			}
		}
		return err
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

func (r *ProductRepository) AddVariantOptionValue(ctx context.Context, variantID uuid.UUID, valueID uuid.UUID) error {
	query := `
		INSERT INTO variant_option_values
		(variant_id, option_value_id)
		VALUES ($1, $2)
	`
	_, err := r.execFromCtx(ctx).ExecContext(ctx, query, variantID, valueID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation -> invalid reference
				return product.ErrInvalidVariant
			case "23502": // not_null_violation
				return product.ErrInvalidVariantOptValInput
			}
		}
		return err
	}

	return nil
}

func (r *ProductRepository) GetVariantByID(ctx context.Context, id uuid.UUID) (*product.Variant, error) {
	query := `
		SELECT id, product_id, sku, price, stock, image_url, options 
		FROM variants
		WHERE id = $1
	`

	var v product.Variant
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &v, query, id); err != nil {
		return nil, fmt.Errorf("get variant by id: %w", err)
	}

	return &v, nil
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidVariant
			case "23502":
				return product.ErrInvalidVariantInput
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return product.ErrVariantNotFound
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidVariant
			case "23502":
				return product.ErrInvalidVariantInput
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return product.ErrVariantNotFound
	}

	return nil
}

func (r *ProductRepository) ListVariantsByProductID(ctx context.Context, productID uuid.UUID) ([]product.Variant, error) {
	query := `
		SELECT id, product_id, sku, price, stock, image_url
		FROM variants
		WHERE product_id = $1
	`

	var variants []product.Variant
	if err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &variants, query, productID); err != nil {
		return nil, fmt.Errorf("list variants by product id: %w", err)
	}

	return variants, nil
}

func (r *ProductRepository) RemoveVariant(ctx context.Context, variantID uuid.UUID) error {
	query := `
		DELETE FROM variants 
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, variantID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidVariant
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify variant deletion: %w", err)
	}

	if rows == 0 {
		return product.ErrVariantNotFound
	}

	return nil
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	query := `
		INSERT INTO products (
			store_id, name, description, category
		) VALUES (
			:store_id, :name, :description, :category
		)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, p)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				return product.ErrProductInvalidStore
			case "23505": // unique_violation
				return product.ErrProductAlreadyExists
			case "23502": // not_null_violation
				return product.ErrInvalidProductInput
			}
		}
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.ID); err != nil {
			return fmt.Errorf("scan product id: %w", err)
		}
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
			if err := r.AddOptionValue(ctx, p.ID, optionID, val.Value); err != nil {
				return err
			}
		}
	}

	// 4. Insert variants if HasVariants == true
	if p.HasVariants {
		for _, v := range p.Variants {
			// mapping - mv top level later
			vi := &product.Variant{
				ProductID: p.ID,
				SKU:       v.SKU,
				Price:     v.Price,
				Stock:     v.Stock,
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

func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	query := `
		SELECT
			p.id,
			p.store_id,
			p.name,
			p.description,
			p.category,
			EXISTS (
				SELECT 1 FROM variants v WHERE v.product_id = p.id
			) AS has_variants
		FROM products p
		WHERE p.id = $1;
	`

	var p product.Product
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &p, query, id); err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &p, nil
}

func (r *ProductRepository) UpdateProductStock(ctx context.Context, productID uuid.UUID, stock int) error {
	params := map[string]interface{}{
		"product_id": productID,
		"stock":      stock,
	}

	query := `
		UPDATE products
		SET stock = :stock,
				updated_at = NOW()
		WHERE id = :product_id
	`
	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidProduct
			case "23502":
				return product.ErrInvalidProductInput
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return product.ErrProductNotFound
	}

	return nil

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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidProduct
			case "23502":
				return product.ErrInvalidProductInput
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) ListProductsByStore(ctx context.Context, storeID uuid.UUID) ([]product.ProductListItem, error) {
	query := `
	SELECT
		p.id,
		p.store_id,
		p.name,
		p.description,
		p.category,
		p.created_at,
		p.updated_at,

		img.url        AS image_url,
		img.is_primary AS is_primary,

		EXISTS (
			SELECT 1 FROM variants v WHERE v.product_id = p.id
		) AS has_variants,

		COUNT(v.id)  AS variant_count,
		MIN(v.price) AS min_price,
		MAX(v.price) AS max_price

	FROM products p

	LEFT JOIN LATERAL (
		SELECT url, is_primary
		FROM product_images pi
		WHERE pi.product_id = p.id
		ORDER BY pi.is_primary DESC, pi.position ASC
		LIMIT 1
	) img ON true

	LEFT JOIN variants v ON v.product_id = p.id
	WHERE p.store_id = $1
	GROUP BY p.id, img.url, img.is_primary;
	`

	var products []product.ProductListItem
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &products, query, storeID)
	return products, err
}

func (r *ProductRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	query := `
		DELETE FROM products 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, productID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return product.ErrInvalidProduct
			}
		}
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify product deletion: %w", err)
	}

	if rows == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) IsOptionUsed(ctx context.Context, optionID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM variant_option_values vov
			JOIN product_option_values pov
				ON vov.option_value_id = pov.id
			WHERE pov.product_option_id = $1
		)
	`
	var exists bool
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &exists, query, optionID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ProductRepository) IsOptionValueUsed(ctx context.Context, optionValueID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM variant_option_values
			WHERE option_value_id = $1
		)
	`
	var exists bool
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &exists, query, optionValueID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ProductRepository) GetFullProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	// 1. Fetch base product info
	p, err := r.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get base product: %w", err)
	}

	// 2. Fetch images
	queryImages := `
        SELECT url
        FROM product_images
        WHERE product_id = $1
        ORDER BY is_primary DESC, position ASC
    `
	var images []string
	if err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &images, queryImages, id); err != nil {
		return nil, fmt.Errorf("list product images: %w", err)
	}
	p.Images = images

	// 3. Fetch options and their values
	options, err := r.ListOptionsByProductID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list options: %w", err)
	}

	for i, opt := range options {
		values, err := r.ListOptionValuesByOptionID(ctx, opt.ID)
		if err != nil {
			return nil, fmt.Errorf("list option values: %w", err)
		}
		options[i].Values = values
	}
	p.Options = options

	// 4. Fetch variants
	variants, err := r.ListVariantsByProductID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list variants: %w", err)
	}

	// 5. Attach option values to each variant (dictionary)
	for i := range variants {
		variant := &variants[i]
		queryOptVals := `
            SELECT pov.id, po.name AS option_name, pov.value
            FROM variant_option_values vov
            JOIN product_option_values pov ON vov.option_value_id = pov.id
            JOIN product_options po ON pov.product_option_id = po.id
            WHERE vov.variant_id = $1
        `
		rows, err := r.execFromCtx(ctx).QueryContext(ctx, queryOptVals, variant.ID)
		if err != nil {
			return nil, fmt.Errorf("list variant options: %w", err)
		}

		optsMap := make(map[string]string)
		defer rows.Close()
		for rows.Next() {
			var id uuid.UUID
			var name, value string
			if err := rows.Scan(&id, &name, &value); err != nil {
				return nil, fmt.Errorf("scan variant option: %w", err)
			}
			optsMap[name] = value
		}
		variant.Options = optsMap
	}
	p.Variants = variants

	return p, nil
}
