-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS menu (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    menu_name VARCHAR(255) NOT NULL,
    menu_route VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT NULL,
    updated_by VARCHAR(255) DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_menu_route ON menu(menu_route);
-- +goose StatementEnd

