-- Create inventories table
CREATE TABLE inventories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID REFERENCES users(id),
    name TEXT NOT NULL,
    category TEXT,
    stock INTEGER,
    price NUMERIC,
    images TEXT,
    unit TEXT,
    packaging TEXT,
    description TEXT,
    location TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);