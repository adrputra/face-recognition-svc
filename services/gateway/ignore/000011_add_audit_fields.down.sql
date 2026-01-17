-- +goose Down
-- +goose StatementBegin
-- Remove audit fields from institution table
DROP INDEX IF EXISTS idx_institution_deleted_at;
ALTER TABLE institution 
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from parameter table
DROP INDEX IF EXISTS idx_parameter_created_at;
DROP INDEX IF EXISTS idx_parameter_deleted_at;
ALTER TABLE parameter 
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from model_training table
DROP INDEX IF EXISTS idx_model_training_updated_at;
DROP INDEX IF EXISTS idx_model_training_deleted_at;
ALTER TABLE model_training 
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from face_datasets table
DROP INDEX IF EXISTS idx_face_datasets_updated_at;
DROP INDEX IF EXISTS idx_face_datasets_deleted_at;
ALTER TABLE face_datasets 
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from users table
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users 
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from menu_mapping table
DROP INDEX IF EXISTS idx_menu_mapping_deleted_at;
ALTER TABLE menu_mapping 
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from menu table
DROP INDEX IF EXISTS idx_menu_deleted_at;
ALTER TABLE menu 
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;

-- Remove audit fields from role table
DROP INDEX IF EXISTS idx_role_deleted_at;
ALTER TABLE role 
DROP COLUMN IF EXISTS deleted_at,
DROP COLUMN IF EXISTS deleted_by;
-- +goose StatementEnd
