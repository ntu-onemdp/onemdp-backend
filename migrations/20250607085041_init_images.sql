-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.images
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    image bytea NOT NULL,
    PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.images;
-- +goose StatementEnd
