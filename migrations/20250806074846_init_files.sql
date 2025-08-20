-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.files (
    file_id text NOT NULL,
    author text NOT NULL,
    filename text NOT NULL,
    gcs_filename text NOT NULL,
    status text NOT NULL DEFAULT 'available'::text,
    time_created timestamp with time zone NOT NULL DEFAULT now(),
    time_deleted timestamp with time zone,
    deleted_by text,
    file_group text,
    PRIMARY KEY (file_id),
    CONSTRAINT fk_files_author_users_uid FOREIGN KEY (author) REFERENCES public.users (uid) MATCH SIMPLE ON UPDATE CASCADE ON DELETE NO ACTION NOT VALID
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.files;
-- +goose StatementEnd