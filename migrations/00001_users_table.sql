-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email citext NOT NULL UNIQUE,
    password text NOT NULL,
    name text,
    profile_picture text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    version integer NOT NULL DEFAULT 1
);

CREATE OR REPLACE FUNCTION users_set_updated_at_and_increment_version()
RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_at := now();
        NEW.version := COALESCE(OLD.version, 0) + 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at_and_increment_version_trg
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION users_set_updated_at_and_increment_version();

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TRIGGER IF EXISTS users_set_updated_at_and_increment_version_trg ON users;
DROP FUNCTION IF EXISTS users_set_updated_at_and_increment_version();
DROP TABLE IF EXISTS users;

COMMIT;
-- +goose StatementEnd
