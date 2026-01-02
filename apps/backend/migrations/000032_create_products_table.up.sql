-- Main product record
CREATE TABLE products (
    id UUID PRIMARY KEY,
    merchant_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index to quickly find all products for a merchant
CREATE INDEX idx_products_merchant_id ON products(merchant_id);
