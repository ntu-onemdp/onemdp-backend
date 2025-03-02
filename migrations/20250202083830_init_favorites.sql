-- +goose Up
-- +goose StatementBegin
-- Table: public.favorites
CREATE TABLE IF NOT EXISTS public.favorites
(
    username text COLLATE pg_catalog."default" NOT NULL,
    content_id uuid NOT NULL,
    content_type text COLLATE pg_catalog."default" NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    CONSTRAINT favorites_pkey PRIMARY KEY (username, content_id),
    CONSTRAINT fk_favorites_username_users_username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.favorites
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.favorites FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.favorites TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.favorites TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.favorites;
-- +goose StatementEnd
