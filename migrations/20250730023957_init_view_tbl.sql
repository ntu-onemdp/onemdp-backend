-- +goose Up
-- +goose StatementBegin
-- Remove views column from threads and articles
ALTER TABLE IF EXISTS PUBLIC.ARTICLES
DROP COLUMN IF EXISTS VIEWS;

ALTER TABLE IF EXISTS PUBLIC.THREADS
DROP COLUMN IF EXISTS VIEWS;

-- Create new views table
CREATE TABLE
    IF NOT EXISTS PUBLIC.VIEWS (
        uid text NOT NULL,
        content_id text NOT NULL,
        PRIMARY KEY (uid, content_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Restore the columns first
ALTER TABLE IF EXISTS public.articles
ADD COLUMN IF NOT EXISTS views integer DEFAULT 0;

ALTER TABLE IF EXISTS public.threads
ADD COLUMN IF NOT EXISTS views integer DEFAULT 0;

-- Set 'views' on articles: count all views with content_id=<article_id>
UPDATE public.articles
SET
    views = COALESCE(
        (
            SELECT
                COUNT(*)
            FROM
                public.views
            WHERE
                content_id = articles.article_id
        ),
        0
    );

-- Set 'views' on threads: count all views with content_id=<thread_id>
UPDATE public.threads
SET
    views = COALESCE(
        (
            SELECT
                COUNT(*)
            FROM
                public.views
            WHERE
                content_id = threads.thread_id
        ),
        0
    );

-- Now drop the views table
DROP TABLE IF EXISTS public.views;

-- +goose StatementEnd