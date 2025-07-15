-- +goose Up
-- +goose StatementBegin
INSERT INTO public.users (uid, name, email, role, semester)
SELECT uuid_generate_v4(), 'NEXTJS_PROXY', 'NEXTJS@NEXTJS.ADMIN', 'student', 'NA'
WHERE NOT EXISTS (
    SELECT 1 FROM public.users WHERE name = 'NEXTJS_PROXY'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.users
WHERE name = 'NEXTJS_PROXY'
  AND email = 'NEXTJS@NEXTJS.ADMIN';
-- +goose StatementEnd
