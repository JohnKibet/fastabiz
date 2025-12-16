-- Add currency column to payments
ALTER TABLE payments
ADD COLUMN currency CHAR(3) NOT NULL DEFAULT 'KES',
ADD CONSTRAINT currency_code_check CHECK (char_length(currency) = 3);

-- Change payments.amount to BIGINT cents
ALTER TABLE payments
ALTER COLUMN amount TYPE BIGINT
USING (amount * 100)::BIGINT;