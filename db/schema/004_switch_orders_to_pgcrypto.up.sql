-- Enable pgcrypto (if missing) and switch orders.id default to gen_random_uuid()
-- Safe to run: only alters the column default, does not modify existing rows

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'id'
    ) THEN
        ALTER TABLE public.orders ALTER COLUMN id SET DEFAULT gen_random_uuid();
    END IF;
END$$;
