-- +goose Up
-- +goose StatementBegin
-- Table: public.auth

CREATE TABLE IF NOT EXISTS public.auth
(
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    role text COLLATE pg_catalog."default" NOT NULL DEFAULT 'student'::text,
    CONSTRAINT auth_pkey PRIMARY KEY (username),
    CONSTRAINT fk_auth_role_permissions_role FOREIGN KEY (role)
        REFERENCES public.permissions (role) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION,
    CONSTRAINT fk_auth_username_users_username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.auth
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.auth FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.auth TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.auth TO onemdp_db_rw_dev;
-- Index: fki_fk_auth_role_permissions_role

CREATE INDEX IF NOT EXISTS fki_fk_auth_role_permissions_role
    ON public.auth USING btree
    (role COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: fki_fk_auth_username_users_username

CREATE INDEX IF NOT EXISTS fki_fk_auth_username_users_username
    ON public.auth USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.auth;
DROP INDEX IF EXISTS public.fki_fk_auth_role_permissions_role;
DROP INDEX IF EXISTS public.fki_fk_auth_username_users_username;

-- +goose StatementEnd
