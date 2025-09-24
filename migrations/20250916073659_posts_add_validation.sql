-- +goose Up
-- +goose StatementBegin
CREATE TYPE VALIDATION_STATUS AS ENUM('unverified', 'validated', 'refuted');
ALTER TABLE IF EXISTS public.posts
ADD COLUMN validation_status VALIDATION_STATUS NOT NULL DEFAULT 'unverified';
ALTER TABLE IF EXISTS public.posts
ADD COLUMN validated_by TEXT;
ALTER TABLE IF EXISTS public.posts
ADD CONSTRAINT fk_posts_validated_by_users_uid FOREIGN KEY (validated_by) REFERENCES public.users (uid) MATCH SIMPLE ON UPDATE CASCADE ON DELETE
SET NULL NOT VALID;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.posts DROP COLUMN validation_status;
ALTER TABLE IF EXISTS public.posts DROP COLUMN validated_by;
DROP TYPE IF EXISTS VALIDATION_STATUS;
ALTER TABLE IF EXISTS public.posts DROP CONSTRAINT fk_posts_validated_by_users_uid;
-- +goose StatementEnd