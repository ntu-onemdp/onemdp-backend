-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.articles
    ADD COLUMN preview text DEFAULT '' NOT NULL;

UPDATE public.articles
    SET preview = LEFT(content, 100);

ALTER TABLE public.articles
    ALTER COLUMN preview DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.articles
    DROP COLUMN IF EXISTS preview;
-- +goose StatementEnd
