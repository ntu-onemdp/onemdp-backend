-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.comments
    ALTER COLUMN comment_id TYPE text;
ALTER TABLE IF EXISTS public.comments
    ALTER COLUMN comment_id DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.comments
    ALTER COLUMN comment_id TYPE uuid DEFAULT uuid_generate_v4();
-- +goose StatementEnd
