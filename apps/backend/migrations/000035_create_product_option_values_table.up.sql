-- Values for each option (e.g., Size → Small, Medium, Large)
CREATE TABLE product_option_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_option_id UUID NOT NULL REFERENCES product_options(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    position INT DEFAULT 0
);

-- Ensure unique values per option
CREATE UNIQUE INDEX ux_option_value
ON product_option_values(product_option_id, LOWER(value));

-- Index for quick lookup of values by option
CREATE INDEX idx_option_values_option_id ON product_option_values(product_option_id);
