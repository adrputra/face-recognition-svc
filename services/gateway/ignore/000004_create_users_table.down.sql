-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_institution_id;
DROP INDEX IF EXISTS idx_role_id;
DROP INDEX IF EXISTS idx_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

