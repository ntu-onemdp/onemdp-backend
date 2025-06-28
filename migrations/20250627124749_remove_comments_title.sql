-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.comments DROP COLUMN IF EXISTS title;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.comments
    ADD COLUMN IF NOT EXISTS title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text;
-- +goose StatementEnd
