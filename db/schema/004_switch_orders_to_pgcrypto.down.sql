-- Revert orders.id default back to uuid_generate_v4()

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'id'
    ) THEN
        ALTER TABLE public.orders ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;
END$$;
