-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_menu_route;
DROP TABLE IF EXISTS menu;
-- +goose StatementEnd

