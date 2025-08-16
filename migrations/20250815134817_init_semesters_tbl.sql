-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.semesters
(
    semester text NOT NULL,
    code text NOT NULL,
    is_current boolean NOT NULL DEFAULT FALSE,
    PRIMARY KEY (semester)
);

-- Ensure that only one row can be true at any time.
CREATE UNIQUE INDEX unique_current_semester ON semesters (is_current)
WHERE is_current;

-- Insert the first semester (hardcoded) into the database
INSERT INTO public.semesters (SEMESTER, CODE, IS_CURRENT) VALUES ('AY2025/01', '1q2w3e', TRUE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.semesters
-- +goose StatementEnd
