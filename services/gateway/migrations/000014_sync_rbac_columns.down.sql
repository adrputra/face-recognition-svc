-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    DROP COLUMN IF EXISTS cover_photo,
    DROP COLUMN IF EXISTS profile_photo;

ALTER TABLE role
    DROP COLUMN IF EXISTS is_administrator;
-- +goose StatementEnd
