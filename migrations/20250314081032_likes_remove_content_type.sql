-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.likes DROP COLUMN IF EXISTS content_type;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.likes
ADD COLUMN IF NOT EXISTS content_type text
-- +goose StatementEnd
