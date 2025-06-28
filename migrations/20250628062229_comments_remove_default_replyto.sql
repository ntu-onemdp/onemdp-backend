-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.comments
    ALTER COLUMN reply_to DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.comments
    ALTER COLUMN reply_to SET DEFAULT '[deleted]';
-- +goose StatementEnd
