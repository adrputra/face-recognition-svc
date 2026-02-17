-- +goose Up
-- +goose StatementBegin
-- Add missing audit fields to role table
ALTER TABLE role 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_role_deleted_at ON role(deleted_at);

-- Add missing audit fields to menu table
ALTER TABLE menu 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_menu_deleted_at ON menu(deleted_at);

-- Add missing audit fields to menu_mapping table
ALTER TABLE menu_mapping 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_menu_mapping_deleted_at ON menu_mapping(deleted_at);

-- Add missing audit fields to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON users(updated_at);

-- Add missing audit fields to face_datasets table
ALTER TABLE face_datasets 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_face_datasets_deleted_at ON face_datasets(deleted_at);
CREATE INDEX IF NOT EXISTS idx_face_datasets_updated_at ON face_datasets(updated_at);

-- Add missing audit fields to model_training table
ALTER TABLE model_training 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_model_training_deleted_at ON model_training(deleted_at);
CREATE INDEX IF NOT EXISTS idx_model_training_updated_at ON model_training(updated_at);

-- Add missing audit fields to parameter table
ALTER TABLE parameter 
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_by VARCHAR(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_parameter_deleted_at ON parameter(deleted_at);
CREATE INDEX IF NOT EXISTS idx_parameter_created_at ON parameter(created_at);

-- Add missing audit fields to institution table
ALTER TABLE institution 
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL,
ADD COLUMN IF NOT EXISTS deleted_by VARCHAR(200) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_institution_deleted_at ON institution(deleted_at);

-- Update existing records to set updated_at = created_at for tables that just got updated_at
UPDATE users SET updated_at = created_at WHERE updated_at IS NULL;
UPDATE face_datasets SET updated_at = created_at WHERE updated_at IS NULL;
UPDATE model_training SET updated_at = created_at WHERE updated_at IS NULL;
UPDATE parameter SET created_at = updated_at WHERE created_at IS NULL;
-- +goose StatementEnd
