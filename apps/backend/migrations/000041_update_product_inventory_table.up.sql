ALTER TABLE product_inventory
ADD COLUMN price numeric(10,2) NOT NULL DEFAULT 0;

ALTER TABLE product_inventory
ADD COLUMN currency char(3) NOT NULL DEFAULT 'USD';

ALTER TABLE product_inventory
ADD CONSTRAINT product_inventory_price_check
CHECK (price >= 0);
