-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.posts
    ADD COLUMN is_header boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.posts
    DROP COLUMN IF EXISTS is_header;
-- +goose StatementEnd
