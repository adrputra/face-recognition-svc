-- +goose Up
-- +goose StatementBegin
ALTER TABLE role
    ADD COLUMN IF NOT EXISTS is_administrator BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE "user"
    ADD COLUMN IF NOT EXISTS profile_photo VARCHAR(500) DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS cover_photo VARCHAR(500) DEFAULT NULL;
-- +goose StatementEnd
