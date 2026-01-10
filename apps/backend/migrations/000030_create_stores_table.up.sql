CREATE EXTENSION IF NOT EXISTS citext;

-- Create stores table
CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    name_normalized CITEXT NOT NULL,
    logo_url TEXT,
    location TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_stores_owner_id ON stores(owner_id);
CREATE UNIQUE INDEX uniq_store_name_per_owner ON stores(owner_id, name_normalized);