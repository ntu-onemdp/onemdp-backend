-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ALTER COLUMN date_created SET DEFAULT now();

ALTER TABLE IF EXISTS public.users
    ALTER COLUMN date_created SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.users
    ALTER COLUMN date_created DROP DEFAULT;

ALTER TABLE IF EXISTS public.users
    ALTER COLUMN date_created DROP NOT NULL;
-- +goose StatementEnd
