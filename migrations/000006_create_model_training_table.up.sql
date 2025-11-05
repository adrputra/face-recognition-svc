-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS model_training (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    institution_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    is_used VARCHAR(50) DEFAULT 'N',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_institution_id ON model_training(institution_id);
CREATE INDEX IF NOT EXISTS idx_status ON model_training(status);
CREATE INDEX IF NOT EXISTS idx_is_used ON model_training(is_used);
CREATE INDEX IF NOT EXISTS idx_created_at ON model_training(created_at);
-- +goose StatementEnd

