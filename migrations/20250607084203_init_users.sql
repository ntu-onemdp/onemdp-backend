-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    uid text,
    name text,
    email text NOT NULL,
    role text NOT NULL DEFAULT 'student',
    date_created timestamp with time zone,
    date_removed timestamp with time zone,
    semester text,
    profile_photo bytea,
    status text NOT NULL DEFAULT 'active',
    karma integer NOT NULL DEFAULT 0,
    PRIMARY KEY (uid)
);

ALTER TABLE IF EXISTS public.users
    OWNER to onemdp_db_admin_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users
-- +goose StatementEnd
