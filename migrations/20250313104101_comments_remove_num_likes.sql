-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.comments DROP COLUMN IF EXISTS num_likes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
LTER TABLE IF EXISTS public.comments
ADD COLUMN IF NOT EXISTS num_likes integer NOT NULL DEFAULT 0;
-- +goose StatementEnd
