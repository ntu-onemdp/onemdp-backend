-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.karma (
    semester text NOT NULL,
    create_thread integer NOT NULL DEFAULT 0,
    create_article integer NOT NULL DEFAULT 0,
    create_comment integer NOT NULL DEFAULT 0,
    create_post integer NOT NULL DEFAULT 0,
    "like" integer NOT NULL DEFAULT 0,
    PRIMARY KEY (semester),
    CONSTRAINT fk_karma_semester_semesters_semester FOREIGN KEY (semester) REFERENCES public.semesters (semester) MATCH SIMPLE ON UPDATE CASCADE ON DELETE CASCADE NOT VALID
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.karma -- +goose StatementEnd