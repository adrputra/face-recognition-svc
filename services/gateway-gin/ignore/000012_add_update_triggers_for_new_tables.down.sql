-- +goose Down
-- +goose StatementBegin
-- Drop triggers
DROP TRIGGER IF EXISTS update_institution_updated_at ON institution;
DROP TRIGGER IF EXISTS update_parameter_updated_at ON parameter;
DROP TRIGGER IF EXISTS update_model_training_updated_at ON model_training;
DROP TRIGGER IF EXISTS update_face_datasets_updated_at ON face_datasets;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- +goose StatementEnd
