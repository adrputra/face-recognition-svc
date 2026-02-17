-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parameter (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    value TEXT NOT NULL,
    description VARCHAR(500) DEFAULT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_updated_at ON parameter(updated_at);
-- +goose StatementEnd

