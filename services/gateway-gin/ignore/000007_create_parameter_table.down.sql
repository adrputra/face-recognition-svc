-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_updated_at;
DROP TABLE IF EXISTS parameter;
-- +goose StatementEnd

