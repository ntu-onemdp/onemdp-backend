-- +goose Up
-- +goose StatementBegin
-- Table: public.articles

CREATE TABLE IF NOT EXISTS public.articles
(
    article_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    author text COLLATE pg_catalog."default" DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL,
    num_likes integer NOT NULL DEFAULT 0,
    num_comments integer NOT NULL DEFAULT 0,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_activity timestamp with time zone NOT NULL DEFAULT now(),
    views integer NOT NULL DEFAULT 1,
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    content text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT articles_pkey PRIMARY KEY (article_id),
    CONSTRAINT fk_articles_author_users_username FOREIGN KEY (author)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.articles
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.articles FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.articles TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.articles TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.articles;
-- +goose StatementEnd
