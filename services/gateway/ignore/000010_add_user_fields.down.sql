-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
DROP COLUMN IF EXISTS cover_photo,
DROP COLUMN IF EXISTS profile_photo,
DROP COLUMN IF EXISTS religion,
DROP COLUMN IF EXISTS gender,
DROP COLUMN IF EXISTS address,
DROP COLUMN IF EXISTS phone_number;
-- +goose StatementEnd

