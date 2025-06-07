-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users ADD CONSTRAINT users_email_key UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.users DROP CONSTRAINT IF EXISTS users_email_key;
-- +goose StatementEnd
