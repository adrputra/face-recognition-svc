-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_menu_role;
DROP INDEX IF EXISTS idx_role_id;
DROP INDEX IF EXISTS idx_menu_id;
DROP TABLE IF EXISTS menu_mapping;
-- +goose StatementEnd

