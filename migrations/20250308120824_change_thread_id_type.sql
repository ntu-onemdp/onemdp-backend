-- This migration changes the id type of threads from UUID to nanoid
-- +goose Up
-- +goose StatementBegin
-- Remove foreign key constraint
ALTER TABLE IF EXISTS public.posts DROP CONSTRAINT IF EXISTS fk_posts_thread_id_threads_pk;

-- Change column data type in threads
ALTER TABLE public.threads
    ALTER COLUMN thread_id TYPE text;
ALTER TABLE IF EXISTS public.threads
    ALTER COLUMN thread_id DROP DEFAULT;

-- Change column data type in posts
ALTER TABLE public.posts
    ALTER COLUMN thread_id TYPE text;

-- Re-add foreign key constraint
ALTER TABLE IF EXISTS public.posts
    ADD CONSTRAINT fk_posts_thread_id_threads_pk FOREIGN KEY (thread_id) REFERENCES threads(thread_id);
    
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint
ALTER TABLE IF EXISTS public.posts DROP CONSTRAINT IF EXISTS fk_posts_thread_id_threads_pk;

-- Revert data type change in posts
ALTER TABLE public.posts
    ALTER COLUMN thread_id TYPE uuid;

-- Revert data type change in threads
ALTER TABLE public.threads
    ALTER COLUMN thread_id TYPE uuid;
ALTER TABLE IF EXISTS public.threads
    ALTER COLUMN thread_id SET DEFAULT gen_random_uuid();

-- Re-add foreign key constraint
ALTER TABLE IF EXISTS public.posts
    ADD CONSTRAINT fk_posts_thread_id_threads_pk FOREIGN KEY (thread_id) REFERENCES threads(thread_id);
-- +goose StatementEnd
