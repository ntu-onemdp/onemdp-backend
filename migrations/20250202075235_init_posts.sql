-- +goose Up
-- +goose StatementBegin
-- Table: public.posts
CREATE TABLE IF NOT EXISTS public.posts
(
    post_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    author text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    thread_id uuid NOT NULL,
    reply_to text COLLATE pg_catalog."default",
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text,
    content text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    num_likes integer NOT NULL DEFAULT 0,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_edited timestamp with time zone NOT NULL DEFAULT now(),
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    CONSTRAINT posts_pkey PRIMARY KEY (post_id),
    CONSTRAINT fk_posts_author_users_username FOREIGN KEY (author)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET DEFAULT,
    CONSTRAINT fk_posts_reply_users_username FOREIGN KEY (reply_to)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_posts_thread_id_threads_pk FOREIGN KEY (thread_id)
        REFERENCES public.threads (thread_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.posts
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.posts FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.posts TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.posts TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.posts;
-- +goose StatementEnd
