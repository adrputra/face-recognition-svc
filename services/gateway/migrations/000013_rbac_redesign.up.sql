-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$
BEGIN
    CREATE TYPE role_scope AS ENUM ('system', 'institution');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    CREATE TYPE feature_type AS ENUM ('menu', 'permission', 'system');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(150) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    short_name VARCHAR(200) DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    profile_photo VARCHAR(500) DEFAULT NULL,
    cover_photo VARCHAR(500) DEFAULT NULL,
    CONSTRAINT uq_user_username UNIQUE (username),
    CONSTRAINT uq_user_email UNIQUE (email),
    CONSTRAINT chk_user_is_active CHECK (is_active IN (TRUE, FALSE))
);

CREATE TABLE IF NOT EXISTS institution (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(80) NOT NULL UNIQUE,
    address VARCHAR(255) DEFAULT NULL,
    phone_number VARCHAR(80) DEFAULT NULL,
    email VARCHAR(255) DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_institution (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    institution_id UUID NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_institutions_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_institutions_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_user_institutions_status CHECK (status IN ('active', 'suspended', 'invited', 'left')),
    CONSTRAINT uq_user_institutions UNIQUE (user_id, institution_id)
);

CREATE TABLE IF NOT EXISTS role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    description VARCHAR(500) DEFAULT NULL,
    scope role_scope NOT NULL DEFAULT 'institution',
    institution_id UUID DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_administrator BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_roles_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_role_scope_institution CHECK (
        (scope = 'system' AND institution_id IS NULL) OR
        (scope = 'institution' AND institution_id IS NOT NULL)
    ),
    CONSTRAINT uq_roles_scope_name UNIQUE (scope, institution_id, name)
);

CREATE TABLE IF NOT EXISTS user_role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    institution_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_by UUID DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_roles_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_roles_membership FOREIGN KEY (user_id, institution_id)
        REFERENCES user_institution(user_id, institution_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uq_user_roles UNIQUE (user_id, institution_id, role_id)
);

CREATE TABLE IF NOT EXISTS permission (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    service VARCHAR(80) NOT NULL,
    resource VARCHAR(120) NOT NULL,
    action VARCHAR(120) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_high_risk BOOLEAN NOT NULL DEFAULT FALSE,
    description VARCHAR(500) DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_permissions_name UNIQUE (name),
    CONSTRAINT uq_permissions_triplet UNIQUE (service, resource, action),
    CONSTRAINT chk_permission_name_format CHECK (
        name ~ '^[a-z][a-z0-9]*\\.[a-z][a-z0-9_]*\\.[a-z][a-z0-9_]*$'
    )
);

CREATE TABLE IF NOT EXISTS role_permission (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permission(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uq_role_permissions UNIQUE (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS feature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feature_key VARCHAR(120) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(500) DEFAULT NULL,
    feature_type feature_type NOT NULL DEFAULT 'system',
    default_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS institution_feature (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID NOT NULL,
    feature_key VARCHAR(120) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_institution_features_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_institution_features_feature FOREIGN KEY (feature_key) REFERENCES feature(feature_key) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uq_institution_features UNIQUE (institution_id, feature_key)
);

CREATE TABLE IF NOT EXISTS menu (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_key VARCHAR(120) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    route VARCHAR(300) DEFAULT NULL,
    icon VARCHAR(120) DEFAULT NULL,
    parent_id UUID DEFAULT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    feature_key VARCHAR(120) DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_menus_parent FOREIGN KEY (parent_id) REFERENCES menu(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_menus_feature FOREIGN KEY (feature_key) REFERENCES feature(feature_key) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS role_menu (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL,
    menu_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_menus_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_role_menus_menu FOREIGN KEY (menu_id) REFERENCES menu(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uq_role_menus UNIQUE (role_id, menu_id)
);

CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID DEFAULT NULL,
    institution_id UUID DEFAULT NULL,
    permission_name VARCHAR(200) DEFAULT NULL,
    action VARCHAR(120) NOT NULL,
    entity_type VARCHAR(120) NOT NULL,
    entity_id VARCHAR(200) DEFAULT NULL,
    request_id VARCHAR(120) DEFAULT NULL,
    ip_address VARCHAR(80) DEFAULT NULL,
    user_agent VARCHAR(300) DEFAULT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_audit_logs_user FOREIGN KEY (actor_user_id) REFERENCES "user"(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_audit_logs_institution FOREIGN KEY (institution_id) REFERENCES institution(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_institutions_user ON user_institution(user_id);
CREATE INDEX IF NOT EXISTS idx_user_institutions_institution ON user_institution(institution_id);
CREATE INDEX IF NOT EXISTS idx_roles_institution ON role(institution_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_institution_user ON user_role(institution_id, user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_role(role_id);
CREATE INDEX IF NOT EXISTS idx_permissions_service ON permission(service);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permission(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission ON role_permission(permission_id);
CREATE INDEX IF NOT EXISTS idx_menus_parent ON menu(parent_id);
CREATE INDEX IF NOT EXISTS idx_role_menus_role ON role_menu(role_id);
CREATE INDEX IF NOT EXISTS idx_role_menus_menu ON role_menu(menu_id);
CREATE INDEX IF NOT EXISTS idx_institution_features_institution ON institution_feature(institution_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_institution ON audit_log(institution_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_log(actor_user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_log(action);

CREATE OR REPLACE FUNCTION prevent_permission_immutable_update()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.name <> OLD.name
        OR NEW.service <> OLD.service
        OR NEW.resource <> OLD.resource
        OR NEW.action <> OLD.action THEN
        RAISE EXCEPTION 'Permission identity fields are immutable';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_permission_immutability
    BEFORE UPDATE ON permission
    FOR EACH ROW
    EXECUTE FUNCTION prevent_permission_immutable_update();

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_institution_updated_at
    BEFORE UPDATE ON institution
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_institution_updated_at
    BEFORE UPDATE ON user_institution
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_role_updated_at
    BEFORE UPDATE ON role
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permission_updated_at
    BEFORE UPDATE ON permission
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feature_updated_at
    BEFORE UPDATE ON feature
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_institution_feature_updated_at
    BEFORE UPDATE ON institution_feature
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_updated_at
    BEFORE UPDATE ON menu
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
