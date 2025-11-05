-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS menu_mapping (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    menu_id VARCHAR(255) NOT NULL,
    role_id VARCHAR(255) NOT NULL,
    access_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT NULL,
    updated_by VARCHAR(255) DEFAULT NULL,
    CONSTRAINT fk_menu_mapping_menu FOREIGN KEY (menu_id) REFERENCES menu(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_menu_mapping_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_menu_id ON menu_mapping(menu_id);
CREATE INDEX IF NOT EXISTS idx_role_id ON menu_mapping(role_id);
CREATE INDEX IF NOT EXISTS idx_menu_role ON menu_mapping(menu_id, role_id);
-- +goose StatementEnd

