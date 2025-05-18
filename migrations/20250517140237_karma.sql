-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ADD COLUMN karma integer NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    DROP COLUMN IF EXISTS karma;
-- +goose StatementEnd
