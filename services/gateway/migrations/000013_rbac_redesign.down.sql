-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_menu_updated_at ON menu;
DROP TRIGGER IF EXISTS update_institution_feature_updated_at ON institution_feature;
DROP TRIGGER IF EXISTS update_feature_updated_at ON feature;
DROP TRIGGER IF EXISTS update_permission_updated_at ON permission;
DROP TRIGGER IF EXISTS update_role_updated_at ON role;
DROP TRIGGER IF EXISTS update_user_institution_updated_at ON user_institution;
DROP TRIGGER IF EXISTS update_institution_updated_at ON institution;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TRIGGER IF EXISTS enforce_permission_immutability ON permission;
DROP FUNCTION IF EXISTS prevent_permission_immutable_update();

DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS role_menu;
DROP TABLE IF EXISTS menu;
DROP TABLE IF EXISTS institution_feature;
DROP TABLE IF EXISTS feature;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS permission;
DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS role;
DROP TABLE IF EXISTS user_institution;
DROP TABLE IF EXISTS institution;
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS feature_type;
DROP TYPE IF EXISTS role_scope;
-- +goose StatementEnd
