ALTER TABLE users
DROP CONSTRAINT users_role_check;

-- Update role constraint to include 'merchant'
ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('admin', 'driver', 'customer', 'merchant'));
