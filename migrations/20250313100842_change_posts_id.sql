-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.posts
    ALTER COLUMN post_id TYPE text;
ALTER TABLE IF EXISTS public.posts
    ALTER COLUMN post_id DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.posts
    ALTER COLUMN post_id TYPE uuid DEFAULT uuid_generate_v4();
-- +goose StatementEnd
