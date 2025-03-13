-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.favorites
    ALTER COLUMN content_id TYPE text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.favorites
    ALTER COLUMN content_id TYPE uuid;
-- +goose StatementEnd
