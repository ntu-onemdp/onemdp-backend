-- +goose Up
-- +goose StatementBegin
-- Table: public.users

CREATE TABLE IF NOT EXISTS public.users
(
    username text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL DEFAULT '[deleted]'::text,
    date_created timestamp with time zone NOT NULL DEFAULT now(),
    date_removed timestamp with time zone,
    semester integer NOT NULL DEFAULT '-1'::integer,
    password_changed boolean NOT NULL DEFAULT false,
    profile_photo text COLLATE pg_catalog."default",
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT users_pkey PRIMARY KEY (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.permissions FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.permissions TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.permissions TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
