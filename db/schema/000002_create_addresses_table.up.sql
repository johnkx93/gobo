-- ==============================================
-- ADDRESSES TABLE
-- ==============================================

CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address VARCHAR(50) NOT NULL,
    floor VARCHAR(10) NOT NULL,
    unit_no VARCHAR(10) NOT NULL,
    block_tower VARCHAR(25),
    company_name VARCHAR(25),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_addresses_user_id ON addresses(user_id);

-- Add trigger for auto-updating updated_at
CREATE TRIGGER trigger_update_addresses_updated_at
    BEFORE UPDATE ON addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
