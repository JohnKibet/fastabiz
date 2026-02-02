-- Drop existing admin_id foreign key constraint
ALTER TABLE orders
DROP CONSTRAINT IF EXISTS fk_orders_admin_id_users_id;

-- Make admin_id nullable (optional)
ALTER TABLE orders
ALTER COLUMN admin_id DROP NOT NULL;

-- Add merchant_id column
ALTER TABLE orders
ADD COLUMN merchant_id UUID NOT NULL REFERENCES users(id);

ALTER TABLE orders
ADD COLUMN store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
ADD COLUMN product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
ADD COLUMN variant_id UUID NULL REFERENCES variants(id) ON DELETE SET NULL,
ADD COLUMN unit_price NUMERIC(10,2) NOT NULL,
ADD COLUMN currency CHAR(3) NOT NULL DEFAULT 'USD',
ADD COLUMN total NUMERIC(10,2) NOT NULL,
ADD COLUMN product_name TEXT NOT NULL,
ADD COLUMN variant_name TEXT,
ADD COLUMN image_url TEXT;