-- +goose Up
-- +goose StatementBegin
-- Table: public.permissions

CREATE TABLE IF NOT EXISTS public.permissions
(
    role text COLLATE pg_catalog."default" NOT NULL,
    manage_students boolean NOT NULL DEFAULT false,
    manage_staff boolean NOT NULL DEFAULT false,
    manage_roles boolean NOT NULL DEFAULT false,
    manage_posts boolean NOT NULL DEFAULT false,
    CONSTRAINT permissions_pkey PRIMARY KEY (role)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.permissions
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.permissions FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.permissions TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.permissions TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.permissions;
-- +goose StatementEnd
