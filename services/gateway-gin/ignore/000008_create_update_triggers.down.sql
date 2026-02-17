-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_menu_mapping_updated_at ON menu_mapping;
DROP TRIGGER IF EXISTS update_menu_updated_at ON menu;
DROP TRIGGER IF EXISTS update_role_updated_at ON role;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd

