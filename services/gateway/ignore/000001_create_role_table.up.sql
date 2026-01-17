-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    role_name VARCHAR(255) NOT NULL,
    role_desc VARCHAR(500) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT NULL,
    updated_by VARCHAR(255) DEFAULT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_administrator BOOLEAN DEFAULT FALSE,
    institution_id VARCHAR(255) NOT NULL,
    CONSTRAINT fk_role_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_role_name ON role(role_name);
CREATE INDEX IF NOT EXISTS idx_is_active ON role(is_active);
-- +goose StatementEnd
