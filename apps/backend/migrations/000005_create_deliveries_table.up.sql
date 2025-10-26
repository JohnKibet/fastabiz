-- Create deliveries table
CREATE TABLE deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    driver_id   UUID REFERENCES drivers(id) ON DELETE SET NULL,
    status  TEXT NOT NULL CHECK (status IN ('assigned', 'picked_up', 'delivered', 'failed')),
    assigned_at TIMESTAMPTZ DEFAULT now(),
    picked_up_at    TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ
);
