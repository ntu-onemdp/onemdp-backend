-- This migration changes the id type of articles from UUID to nanoid
-- +goose Up
-- +goose StatementBegin
-- Remove foreign key constraint
ALTER TABLE IF EXISTS public.comments DROP CONSTRAINT IF EXISTS fk_article_id_articles_pk;

-- Change column data type in articles
ALTER TABLE public.articles
    ALTER COLUMN article_id TYPE text;
ALTER TABLE IF EXISTS public.articles
    ALTER COLUMN article_id DROP DEFAULT;

-- Change column data type in comments
ALTER TABLE public.comments
    ALTER COLUMN article_id TYPE text;

-- Re-add foreign key constraint
ALTER TABLE IF EXISTS public.comments
    ADD CONSTRAINT fk_article_id_articles_pk FOREIGN KEY (article_id) REFERENCES articles(article_id);
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint
ALTER TABLE IF EXISTS public.comments DROP CONSTRAINT IF EXISTS fk_article_id_articles_pk;

-- Revert column data type in comments
ALTER TABLE public.comments
    ALTER COLUMN article_id TYPE uuid;

-- Revert column data type in articles
ALTER TABLE public.articles
    ALTER COLUMN article_id TYPE uuid;
ALTER TABLE IF EXISTS public.articles
    ALTER COLUMN article_id SET DEFAULT gen_random_uuid();

-- Re-add foreign key constraint
ALTER TABLE IF EXISTS public.comments
    ADD CONSTRAINT fk_article_id_articles_pk FOREIGN KEY (article_id) REFERENCES articles(article_id);
-- +goose StatementEnd
