-- +goose Up
-- +goose StatementBegin
-- Articles table
ALTER TABLE public.articles
DROP CONSTRAINT fk_articles_author_users_uid;

ALTER TABLE public.articles
ADD CONSTRAINT fk_articles_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET DEFAULT
NOT VALID;

-- Comments table
ALTER TABLE public.comments
DROP CONSTRAINT fk_articles_reply_users_uid;

ALTER TABLE public.comments
ADD CONSTRAINT fk_articles_reply_users_uid
FOREIGN KEY (reply_to)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET DEFAULT
NOT VALID;

-- Threads table
ALTER TABLE public.threads
DROP CONSTRAINT fk_theads_author_users_uid;

ALTER TABLE public.threads
ADD CONSTRAINT fk_theads_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET DEFAULT
NOT VALID;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.articles
DROP CONSTRAINT fk_articles_author_users_uid;

ALTER TABLE public.articles
ADD CONSTRAINT fk_articles_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET NULL
NOT VALID;

-- Comments table
ALTER TABLE public.comments
DROP CONSTRAINT fk_articles_reply_users_uid;

ALTER TABLE public.comments
ADD CONSTRAINT fk_articles_reply_users_uid
FOREIGN KEY (reply_to)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET NULL
NOT VALID;

-- Threads table
ALTER TABLE public.threads
DROP CONSTRAINT fk_theads_author_users_uid;

ALTER TABLE public.threads
ADD CONSTRAINT fk_theads_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET NULL
NOT VALID;
-- +goose StatementEnd
