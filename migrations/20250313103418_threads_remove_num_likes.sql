-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.threads DROP COLUMN IF EXISTS num_likes;

ALTER TABLE IF EXISTS public.threads DROP COLUMN IF EXISTS num_replies;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.threads 
ADD COLUMN IF NOT EXISTS num_likes integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.threads 
ADD COLUMN IF NOT EXISTS num_replies integer NOT NULL DEFAULT 0;
-- +goose StatementEnd
