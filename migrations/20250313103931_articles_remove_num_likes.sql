-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.articles DROP COLUMN IF EXISTS num_likes;

ALTER TABLE IF EXISTS public.articles DROP COLUMN IF EXISTS num_comments;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.articles 
ADD COLUMN IF NOT EXISTS num_likes integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.articles 
ADD COLUMN IF NOT EXISTS num_comments integer NOT NULL DEFAULT 0;
-- +goose StatementEnd
