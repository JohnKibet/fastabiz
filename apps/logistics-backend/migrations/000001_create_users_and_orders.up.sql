-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'driver', 'customer')),
    created_at TIMESTAMPZ DEFAULT now(),
    updated_at TIMESTAMPZ DEFAULT now()
);


-- Create orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pickup_address TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'assigned', 'in-transit', 'delivered', 'cancelled')),
    eta TIMESTAMPZ,
    delivery_proof TEXT,
    created_at TIMESTAMPZ DEFAULT now(),
    updated_at TIMESTAMPZ DEFAULT now()
);
