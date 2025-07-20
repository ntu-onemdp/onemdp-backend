-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.threads
ADD COLUMN is_anon boolean NOT NULL DEFAULT false;

ALTER TABLE IF EXISTS public.posts
ADD COLUMN is_anon boolean NOT NULL DEFAULT false;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.threads
DROP COLUMN is_anon;

ALTER TABLE IF EXISTS public.posts
DROP COLUMN is_anon;

-- +goose StatementEnd