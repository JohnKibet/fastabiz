-- Add quantity column to orders
ALTER TABLE orders
ADD COLUMN quantity INTEGER NOT NULL CHECK (quantity > 0);
