-- Migration: fix currency column in payments

BEGIN;

-- ===== Payments =====

-- 1. Ensure currency is VARCHAR(3) instead of CHAR(3)
ALTER TABLE payments
ALTER COLUMN currency TYPE VARCHAR(3);

-- 2. Drop old constraints if they exist
ALTER TABLE payments
DROP CONSTRAINT IF EXISTS currency_code_check;
ALTER TABLE payments
DROP CONSTRAINT IF EXISTS currency_allowed_check;

-- 3. Add stricter regex check
ALTER TABLE payments
ADD CONSTRAINT currency_code_check
CHECK (currency ~ '^[A-Z]{3}$');

-- 4. Optional: whitelist of supported currencies
ALTER TABLE payments
ADD CONSTRAINT currency_allowed_check
CHECK (currency IN ('KES', 'USD', 'EUR', 'GBP'));

COMMIT;
