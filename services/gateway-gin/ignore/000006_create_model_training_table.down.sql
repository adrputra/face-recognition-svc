-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_created_at;
DROP INDEX IF EXISTS idx_is_used;
DROP INDEX IF EXISTS idx_status;
DROP INDEX IF EXISTS idx_institution_id;
DROP TABLE IF EXISTS model_training;
-- +goose StatementEnd

