-- +goose Up
-- +goose StatementBegin
-- Table: public.threads

CREATE TABLE IF NOT EXISTS public.threads
(
    thread_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    author text COLLATE pg_catalog."default" DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    num_likes integer NOT NULL DEFAULT 0,
    num_replies integer NOT NULL DEFAULT 0,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_activity timestamp with time zone NOT NULL DEFAULT now(),
    views integer NOT NULL DEFAULT 1,
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    preview text COLLATE pg_catalog."default",
    CONSTRAINT threads_pkey PRIMARY KEY (thread_id),
    CONSTRAINT fk_theads_author_users_username FOREIGN KEY (author)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.threads
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.threads FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.threads TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.threads TO onemdp_db_rw_dev;


-- Index: fki_fk_theads_author_users_username
CREATE INDEX IF NOT EXISTS fki_fk_theads_author_users_username
    ON public.threads USING btree
    (author COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.threads;

DROP INDEX IF EXISTS public.fki_fk_theads_author_users_username;
-- +goose StatementEnd
