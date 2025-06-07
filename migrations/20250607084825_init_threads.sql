-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.threads
(
    thread_id text NOT NULL,
    author text COLLATE pg_catalog."default" DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_activity timestamp with time zone NOT NULL DEFAULT now(),
    views integer NOT NULL DEFAULT 1,
    flagged boolean NOT NULL DEFAULT false,
    is_available boolean NOT NULL DEFAULT true,
    preview text COLLATE pg_catalog."default",
    CONSTRAINT threads_pkey PRIMARY KEY (thread_id),
    CONSTRAINT fk_theads_author_users_uid FOREIGN KEY (author)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

ALTER TABLE IF EXISTS public.threads
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.threads FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.threads TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.threads TO onemdp_db_rw_dev;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.threads;
-- +goose StatementEnd
