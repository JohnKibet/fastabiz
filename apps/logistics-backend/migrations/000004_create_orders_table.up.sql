-- Create orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pickup_address TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'assigned', 'in-transit', 'delivered', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
