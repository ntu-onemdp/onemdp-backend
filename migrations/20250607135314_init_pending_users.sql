-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.pending_users
(
    email text NOT NULL,
    semester text NOT NULL,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    role text NOT NULL DEFAULT 'student',
    PRIMARY KEY (email)
);

ALTER TABLE IF EXISTS public.pending_users
    OWNER to onemdp_db_admin_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.pending_users
-- +goose StatementEnd
