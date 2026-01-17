-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bucket;
DROP INDEX IF EXISTS idx_dataset;
DROP INDEX IF EXISTS idx_username;
DROP TABLE IF EXISTS face_datasets;
-- +goose StatementEnd

