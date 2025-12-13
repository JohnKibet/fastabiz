-- Used only if product has no variants
CREATE TABLE product_inventory (
    product_id UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    stock INT DEFAULT 0 CHECK (stock >= 0),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
