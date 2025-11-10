-- Add default_address_id column to users table
ALTER TABLE users
ADD COLUMN default_address_id UUID REFERENCES addresses(id) ON DELETE SET NULL;

-- Create index for faster lookups
CREATE INDEX idx_users_default_address_id ON users(default_address_id);
