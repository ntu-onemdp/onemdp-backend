-- +goose Up
-- +goose StatementBegin
-- 1. Add special deleted user (e.g., UID = 0)
INSERT INTO public.users (uid, name, email, role) VALUES ('[deleted]', '[deleted user]', 'N.A.', 'deleted') ON CONFLICT (uid) DO NOTHING;;

-- 2. Set default for 'author' column in files
ALTER TABLE public.files ALTER COLUMN author SET DEFAULT 'deleted';

-- 3. Drop existing foreign key constraint
ALTER TABLE public.files
DROP CONSTRAINT fk_files_author_users_uid;

-- 4. Add new FK constraint
ALTER TABLE public.files
ADD CONSTRAINT fk_files_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE SET DEFAULT
NOT VALID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- 1. Drop the newer FK constraint
ALTER TABLE public.files DROP CONSTRAINT IF EXISTS fk_files_author_users_uid;

-- 2. Remove the default value for 'author'
ALTER TABLE public.files ALTER COLUMN author DROP DEFAULT;

-- 3. Restore the previous FK constraint
ALTER TABLE public.files
ADD CONSTRAINT fk_files_author_users_uid
FOREIGN KEY (author)
REFERENCES public.users(uid)
ON UPDATE CASCADE
ON DELETE NO ACTION
NOT VALID;

-- 4. Delete the "[deleted user]" row
DELETE FROM public.users WHERE uid = 0 AND username = '[deleted user]';
-- +goose StatementEnd
