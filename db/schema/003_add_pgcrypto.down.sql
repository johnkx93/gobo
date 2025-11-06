-- Revert users.id default back to uuid_generate_v4() and drop pgcrypto

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'id'
    ) THEN
        ALTER TABLE public.users ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;
END$$;

-- Drop pgcrypto if present. Note: dropping may fail if other objects depend on it.
DROP EXTENSION IF EXISTS pgcrypto;
