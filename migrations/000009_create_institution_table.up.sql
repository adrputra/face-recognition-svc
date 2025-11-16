-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS institution (
    id VARCHAR(200) NOT NULL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) DEFAULT NULL,
    phone_number VARCHAR(200) DEFAULT NULL,
    email VARCHAR(200) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200) DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(200) DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_institution_email ON institution(email);
CREATE INDEX IF NOT EXISTS idx_institution_name ON institution(name);
-- +goose StatementEnd

