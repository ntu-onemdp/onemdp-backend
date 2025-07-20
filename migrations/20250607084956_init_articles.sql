-- +goose Up
-- +goose StatementBegin
-- Table: public.articles

CREATE TABLE IF NOT EXISTS public.articles
(
    article_id text NOT NULL,
    author text COLLATE pg_catalog."default" DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_activity timestamp with time zone NOT NULL DEFAULT now(),
    views integer NOT NULL DEFAULT 1,
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    content text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT articles_pkey PRIMARY KEY (article_id),
    CONSTRAINT fk_articles_author_users_uid FOREIGN KEY (author)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.articles;
-- +goose StatementEnd
