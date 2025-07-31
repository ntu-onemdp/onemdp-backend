-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    public.favorites (
        uid text NOT NULL,
        content_id text NOT NULL,
        "timestamp" timestamp
        with
            time zone DEFAULT now (),
            PRIMARY KEY (uid, content_id),
            CONSTRAINT fk_favorites_uid_users_uid FOREIGN KEY (uid) REFERENCES public.users (uid) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS PUBLIC.favorites
-- +goose StatementEnd