-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN salt
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
salt character varying COLLATE pg_catalog."default" NOT NULL,
-- +goose StatementEnd
