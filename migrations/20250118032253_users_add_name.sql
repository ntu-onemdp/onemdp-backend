-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ADD COLUMN name character varying NOT NULL DEFAULT 'NA'::character varying;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    DROP COLUMN IF EXISTS name;
-- +goose StatementEnd
