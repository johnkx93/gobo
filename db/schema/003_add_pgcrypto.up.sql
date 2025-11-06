-- Enable pgcrypto and switch users.id default to gen_random_uuid()
-- Run this migration after applying existing migrations (creates extension if missing)

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Switch users.id default from uuid_generate_v4() to gen_random_uuid()
DO $$
BEGIN
    -- Only attempt if table and column exist
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'id'
    ) THEN
        ALTER TABLE public.users ALTER COLUMN id SET DEFAULT gen_random_uuid();
    END IF;
END$$;
