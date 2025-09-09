-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS users_name_trgm_idx ON users USING gin (name gin_trgm_ops);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_name_trgm_idx;
-- +goose StatementEnd