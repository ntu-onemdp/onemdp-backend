-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
  r RECORD;
BEGIN
  FOR r IN (
    SELECT tablename
    FROM pg_tables
    WHERE schemaname = 'public'
      AND tablename <> 'goose_db_version'
  ) LOOP
    EXECUTE 'DROP TABLE IF EXISTS "' || r.tablename || '" CASCADE';
  END LOOP;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- This migration cannot be reversed as it drops the public schema and all its contents.
-- +goose StatementEnd
