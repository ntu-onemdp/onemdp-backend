-- +goose Up
-- +goose StatementBegin
-- Table: public.drafts
CREATE TABLE IF NOT EXISTS public.drafts
(
    draft_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    author text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text,
    content text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_edited timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT drafts_pkey PRIMARY KEY (draft_id),
    CONSTRAINT fk_drafts_author_users_username FOREIGN KEY (author)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.drafts
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.drafts FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.drafts TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.drafts TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.drafts;
-- +goose StatementEnd
