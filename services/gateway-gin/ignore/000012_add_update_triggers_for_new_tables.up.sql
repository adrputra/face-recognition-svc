-- +goose Up
-- +goose StatementBegin
-- Trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for face_datasets table
CREATE TRIGGER update_face_datasets_updated_at
    BEFORE UPDATE ON face_datasets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for model_training table
CREATE TRIGGER update_model_training_updated_at
    BEFORE UPDATE ON model_training
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for parameter table
CREATE TRIGGER update_parameter_updated_at
    BEFORE UPDATE ON parameter
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for institution table
CREATE TRIGGER update_institution_updated_at
    BEFORE UPDATE ON institution
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
