-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_institution_name;
DROP INDEX IF EXISTS idx_institution_email;
DROP TABLE IF EXISTS institution;
-- +goose StatementEnd
