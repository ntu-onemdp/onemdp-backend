-- Set semester to nullable, replace is_active boolean to status TEXT
-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ALTER COLUMN semester DROP NOT NULL;

ALTER TABLE public.users
    RENAME is_active TO status;
ALTER TABLE public.users
    ALTER COLUMN status TYPE text;
ALTER TABLE IF EXISTS public.users
    ALTER COLUMN status SET DEFAULT 'active';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ALTER COLUMN semester SET NOT NULL;

	-- Add the new is_active column with a default value of FALSE
ALTER TABLE public.users
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT FALSE;

-- Convert existing status values to boolean
UPDATE public.users
SET is_active = (CASE WHEN status = 'active' THEN TRUE ELSE FALSE END);

-- Remove the status column
ALTER TABLE public.users
    DROP COLUMN status;
-- +goose StatementEnd	
