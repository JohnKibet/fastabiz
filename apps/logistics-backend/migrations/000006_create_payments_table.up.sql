-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount  NUMERIC(10, 2) NOT NULL,
    method TEXT NOT NULL CHECK (method IN ('stripe', 'paypal', 'mobile_money', 'cash_on_delivery')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    paid_at TIMESTAMPTZ DEFAULT now()
);
