-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_is_active;
DROP INDEX IF EXISTS idx_role_name;
DROP TABLE IF EXISTS role;
-- +goose StatementEnd

