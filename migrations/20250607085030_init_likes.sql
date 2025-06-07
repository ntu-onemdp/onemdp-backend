-- +goose Up
-- +goose StatementBegin
-- Table: public.likes
CREATE TABLE IF NOT EXISTS public.likes
(
    username text COLLATE pg_catalog."default" NOT NULL,
    content_id text NOT NULL,
    "timestamp" timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT likes_pkey PRIMARY KEY (username, content_id),
    CONSTRAINT fk_likes_username_users_uid FOREIGN KEY (username)
        REFERENCES public.users (uid) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.likes
    OWNER to onemdp_db_admin_dev;

REVOKE ALL ON TABLE public.likes FROM onemdp_db_rw_dev;

GRANT ALL ON TABLE public.likes TO onemdp_db_admin_dev;

GRANT DELETE, SELECT, INSERT ON TABLE public.likes TO onemdp_db_rw_dev;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.likes;
-- +goose StatementEnd
