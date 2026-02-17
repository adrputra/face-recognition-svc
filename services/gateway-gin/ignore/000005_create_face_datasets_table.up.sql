-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS face_datasets (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    bucket VARCHAR(255) NOT NULL,
    dataset VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_face_datasets_users FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_username ON face_datasets(username);
CREATE INDEX IF NOT EXISTS idx_dataset ON face_datasets(dataset);
CREATE INDEX IF NOT EXISTS idx_bucket ON face_datasets(bucket);
-- +goose StatementEnd

