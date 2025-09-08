-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- Threads
CREATE INDEX IF NOT EXISTS threads_title_trgm_idx ON threads USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS threads_preview_trgm_idx ON threads USING gin (preview gin_trgm_ops);
-- Posts
CREATE INDEX IF NOT EXISTS posts_content_trgm_idx ON posts USING gin (content gin_trgm_ops);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS posts_content_trgm_idx;
DROP INDEX IF EXISTS threads_preview_trgm_idx;
DROP INDEX IF EXISTS threads_title_trgm_idx;
DROP EXTENSION IF EXISTS pg_trgm;
-- +goose StatementEnd