-- Insert into DB permissions for student, staff and admin
-- +goose Up
-- +goose StatementBegin
INSERT INTO public.permissions(
	role, manage_students, manage_staff, manage_roles, manage_posts)
	VALUES ('student', false, false, false, false);

INSERT INTO public.permissions(
	role, manage_students, manage_staff, manage_roles, manage_posts)
	VALUES ('staff', true, false, false, true);

INSERT INTO public.permissions(
	role, manage_students, manage_staff, manage_roles, manage_posts)
	VALUES ('admin', true, true, true, true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.permissions WHERE role IN ('student', 'staff', 'admin');
-- +goose StatementEnd
