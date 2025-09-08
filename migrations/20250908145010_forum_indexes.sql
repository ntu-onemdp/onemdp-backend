-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_posts_content_gin ON posts USING gin(to_tsvector('english', content));
CREATE INDEX idx_threads_title_gin ON threads USING gin(to_tsvector('english', title));
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_posts_content_gin;
DROP INDEX IF EXISTS idx_threads_title_gin;
-- +goose StatementEnd