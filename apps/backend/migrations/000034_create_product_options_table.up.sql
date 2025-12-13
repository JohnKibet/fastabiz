-- e.g., Weight, Size, Color
CREATE TABLE product_options (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    position INT DEFAULT 0
);

-- Unique option names per product (no duplicate "Size")
CREATE UNIQUE INDEX ux_product_option_name
ON product_options(product_id, LOWER(name));

-- Index for quick lookup of options by product
CREATE INDEX idx_product_options_product_id ON product_options(product_id);
