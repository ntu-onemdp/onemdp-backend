-- +goose Up
-- +goose StatementBegin
-- Table: public.drafts
CREATE TABLE IF NOT EXISTS public.drafts
(
    draft_id text NOT NULL,
    author text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT 'Untitled'::text,
    content text COLLATE pg_catalog."default" NOT NULL DEFAULT 'NA'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    last_edited timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT drafts_pkey PRIMARY KEY (draft_id),
    CONSTRAINT fk_drafts_author_users_uid FOREIGN KEY (author)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.drafts;
-- +goose StatementEnd
