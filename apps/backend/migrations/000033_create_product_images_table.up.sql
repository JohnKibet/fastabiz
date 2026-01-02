-- Multiple images per product
CREATE TABLE product_images (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    position INT DEFAULT 0
);

-- Ensure only one primary image per product
CREATE UNIQUE INDEX ux_product_primary_image
ON product_images(product_id)
WHERE is_primary = true;

-- Index for quick lookup of images per product
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
