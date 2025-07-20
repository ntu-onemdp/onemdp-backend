-- +goose Up
-- +goose StatementBegin
-- Table: public.posts
CREATE TABLE IF NOT EXISTS public.posts
(
    post_id text NOT NULL,
    author text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    thread_id text NOT NULL,
    reply_to text COLLATE pg_catalog."default",
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text,
    content text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_edited timestamp with time zone NOT NULL DEFAULT now(),
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    is_header boolean NOT NULL DEFAULT false,
    CONSTRAINT posts_pkey PRIMARY KEY (post_id),
    CONSTRAINT fk_posts_author_users_uid FOREIGN KEY (author)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET DEFAULT,
    CONSTRAINT fk_posts_reply_users_uid FOREIGN KEY (reply_to)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_posts_thread_id_threads_pk FOREIGN KEY (thread_id)
        REFERENCES public.threads (thread_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.posts;
-- +goose StatementEnd
