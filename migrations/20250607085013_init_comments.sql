-- +goose Up
-- +goose StatementBegin
-- Table: public.comments
CREATE TABLE IF NOT EXISTS public.comments
(
    comment_id text NOT NULL,
    author text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    article_id text NOT NULL,
    reply_to text COLLATE pg_catalog."default" DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text,
    content text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_edited timestamp with time zone NOT NULL DEFAULT now(),
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    CONSTRAINT comments_pkey PRIMARY KEY (comment_id),
    CONSTRAINT fk_article_id_articles_pk FOREIGN KEY (article_id)
        REFERENCES public.articles (article_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_articles_author_users_uid FOREIGN KEY (author)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET DEFAULT,
    CONSTRAINT fk_articles_reply_users_uid FOREIGN KEY (reply_to)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

ALTER TABLE IF EXISTS public.comments
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.comments FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.comments TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.comments TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.comments;
-- +goose StatementEnd
