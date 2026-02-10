-- NOTE:
-- - Assumption: permission.service is 'gateway' for all gateway routes.
-- - Assumption: feature_key values are 'user_management', 'rbac_management',
--   'dataset_management', 'institution_management', and 'param_management'.
-- - Assumption: feature_type stays default 'system', default_enabled = true.
-- - Uses schema defaults for created_at/updated_at and gen_random_uuid().

-- Feature: user_management
-- Routes: GET /api/service/user, GET /api/service/user/detail/:id, PUT /api/service/user,
--         DELETE /api/service/user/:id, GET /api/service/user/institutions,
--         POST /api/service/user/profile-photo, POST /api/service/user/cover-photo
INSERT INTO public.feature (id, feature_key, name, description, feature_type, default_enabled)
SELECT gen_random_uuid(), 'user_management', 'User Management', 'Permissions for user endpoints', 'system', TRUE
WHERE NOT EXISTS (SELECT 1 FROM public.feature WHERE feature_key = 'user_management');

WITH perms AS (
  SELECT 'gateway.user.read' AS name, 'gateway' AS service, 'user' AS resource, 'read' AS action, 'Read user resources' AS description UNION ALL
  SELECT 'gateway.user.update', 'gateway', 'user', 'update', 'Update user resources' UNION ALL
  SELECT 'gateway.user.delete', 'gateway', 'user', 'delete', 'Delete user resources' UNION ALL
  SELECT 'gateway.user.upload_profile_photo', 'gateway', 'user', 'upload_profile_photo', 'Upload user profile photo' UNION ALL
  SELECT 'gateway.user.upload_cover_photo', 'gateway', 'user', 'upload_cover_photo', 'Upload user cover photo'
)
INSERT INTO public.permission (id, name, service, resource, action, is_active, is_high_risk, description)
SELECT gen_random_uuid(), perms.name, perms.service, perms.resource, perms.action, TRUE, FALSE, perms.description
FROM perms
WHERE NOT EXISTS (SELECT 1 FROM public.permission p WHERE p.name = perms.name);

-- Feature: rbac_management
-- Routes: GET /api/service/role, POST /api/service/role/create,
--         GET /api/service/role/mapping, POST /api/service/role/mapping/create,
--         PUT /api/service/role/mapping, DELETE /api/service/role/mapping/:id,
--         GET /api/service/role/menu, POST /api/service/role/menu/create,
--         PUT /api/service/role/menu, DELETE /api/service/role/menu/:id
-- Routes: GET /api/service/permission, POST /api/service/permission, PUT /api/service/permission,
--         POST /api/service/permission/assign, GET /api/service/permission/role/:id
-- Routes: GET /api/service/feature, POST /api/service/feature, PUT /api/service/feature,
--         POST /api/service/feature/institution, GET /api/service/feature/institution/:id
-- Assumption: role/mapping endpoints map to role_mapping.* permissions.
INSERT INTO public.feature (id, feature_key, name, description, feature_type, default_enabled)
SELECT gen_random_uuid(), 'rbac_management', 'RBAC Management', 'Permissions for RBAC administration endpoints', 'system', TRUE
WHERE NOT EXISTS (SELECT 1 FROM public.feature WHERE feature_key = 'rbac_management');

WITH perms AS (
  SELECT 'gateway.role.read' AS name, 'gateway' AS service, 'role' AS resource, 'read' AS action, 'Read roles' AS description UNION ALL
  SELECT 'gateway.role.create', 'gateway', 'role', 'create', 'Create roles' UNION ALL
  SELECT 'gateway.role_mapping.read', 'gateway', 'role_mapping', 'read', 'Read role mappings' UNION ALL
  SELECT 'gateway.role_mapping.create', 'gateway', 'role_mapping', 'create', 'Create role mappings' UNION ALL
  SELECT 'gateway.role_mapping.update', 'gateway', 'role_mapping', 'update', 'Update role mappings' UNION ALL
  SELECT 'gateway.role_mapping.delete', 'gateway', 'role_mapping', 'delete', 'Delete role mappings' UNION ALL
  SELECT 'gateway.menu.read', 'gateway', 'menu', 'read', 'Read menus' UNION ALL
  SELECT 'gateway.menu.create', 'gateway', 'menu', 'create', 'Create menus' UNION ALL
  SELECT 'gateway.menu.update', 'gateway', 'menu', 'update', 'Update menus' UNION ALL
  SELECT 'gateway.menu.delete', 'gateway', 'menu', 'delete', 'Delete menus' UNION ALL
  SELECT 'gateway.permission.read', 'gateway', 'permission', 'read', 'Read permissions' UNION ALL
  SELECT 'gateway.permission.create', 'gateway', 'permission', 'create', 'Create permissions' UNION ALL
  SELECT 'gateway.permission.update', 'gateway', 'permission', 'update', 'Update permissions' UNION ALL
  SELECT 'gateway.permission.assign', 'gateway', 'permission', 'assign', 'Assign role permissions' UNION ALL
  SELECT 'gateway.feature.read', 'gateway', 'feature', 'read', 'Read features' UNION ALL
  SELECT 'gateway.feature.create', 'gateway', 'feature', 'create', 'Create features' UNION ALL
  SELECT 'gateway.feature.update', 'gateway', 'feature', 'update', 'Update features' UNION ALL
  SELECT 'gateway.feature.set_institution', 'gateway', 'feature', 'set_institution', 'Set institution features'
)
INSERT INTO public.permission (id, name, service, resource, action, is_active, is_high_risk, description)
SELECT gen_random_uuid(), perms.name, perms.service, perms.resource, perms.action, TRUE, FALSE, perms.description
FROM perms
WHERE NOT EXISTS (SELECT 1 FROM public.permission p WHERE p.name = perms.name);

-- Feature: dataset_management
-- Routes: GET /api/service/dataset, POST /api/service/dataset, DELETE /api/service/dataset/:id,
--         POST /api/service/dataset/train-model/:id, GET /api/service/dataset/last-train-model/:id,
--         POST /api/service/dataset/model-training-history, GET /api/service/dataset/:institution-id/:id
-- Assumption: dataset routes group under dataset_management.
INSERT INTO public.feature (id, feature_key, name, description, feature_type, default_enabled)
SELECT gen_random_uuid(), 'dataset_management', 'Dataset Management', 'Permissions for dataset endpoints', 'system', TRUE
WHERE NOT EXISTS (SELECT 1 FROM public.feature WHERE feature_key = 'dataset_management');

WITH perms AS (
  SELECT 'gateway.dataset.read' AS name, 'gateway' AS service, 'dataset' AS resource, 'read' AS action, 'Read datasets and training status' AS description UNION ALL
  SELECT 'gateway.dataset.create', 'gateway', 'dataset', 'create', 'Upload/create datasets' UNION ALL
  SELECT 'gateway.dataset.delete', 'gateway', 'dataset', 'delete', 'Delete datasets' UNION ALL
  SELECT 'gateway.dataset.train_model', 'gateway', 'dataset', 'train_model', 'Train dataset model' UNION ALL
  SELECT 'gateway.dataset.model_training_history', 'gateway', 'dataset', 'model_training_history', 'Get model training history'
)
INSERT INTO public.permission (id, name, service, resource, action, is_active, is_high_risk, description)
SELECT gen_random_uuid(), perms.name, perms.service, perms.resource, perms.action, TRUE, FALSE, perms.description
FROM perms
WHERE NOT EXISTS (SELECT 1 FROM public.permission p WHERE p.name = perms.name);

-- Feature: institution_management
-- Routes: GET /api/service/institution, GET /api/service/institution/:id,
--         POST /api/service/institution, PUT /api/service/institution, DELETE /api/service/institution/:id
-- Assumption: institution routes group under institution_management.
INSERT INTO public.feature (id, feature_key, name, description, feature_type, default_enabled)
SELECT gen_random_uuid(), 'institution_management', 'Institution Management', 'Permissions for institution endpoints', 'system', TRUE
WHERE NOT EXISTS (SELECT 1 FROM public.feature WHERE feature_key = 'institution_management');

WITH perms AS (
  SELECT 'gateway.institution.read' AS name, 'gateway' AS service, 'institution' AS resource, 'read' AS action, 'Read institutions' AS description UNION ALL
  SELECT 'gateway.institution.create', 'gateway', 'institution', 'create', 'Create institutions' UNION ALL
  SELECT 'gateway.institution.update', 'gateway', 'institution', 'update', 'Update institutions' UNION ALL
  SELECT 'gateway.institution.delete', 'gateway', 'institution', 'delete', 'Delete institutions'
)
INSERT INTO public.permission (id, name, service, resource, action, is_active, is_high_risk, description)
SELECT gen_random_uuid(), perms.name, perms.service, perms.resource, perms.action, TRUE, FALSE, perms.description
FROM perms
WHERE NOT EXISTS (SELECT 1 FROM public.permission p WHERE p.name = perms.name);

-- Feature: param_management
-- Routes: GET /api/service/param, GET /api/service/param/:id,
--         POST /api/service/param, PUT /api/service/param, DELETE /api/service/param/:id
-- Assumption: parameter routes group under param_management.
INSERT INTO public.feature (id, feature_key, name, description, feature_type, default_enabled)
SELECT gen_random_uuid(), 'param_management', 'Parameter Management', 'Permissions for parameter endpoints', 'system', TRUE
WHERE NOT EXISTS (SELECT 1 FROM public.feature WHERE feature_key = 'param_management');

WITH perms AS (
  SELECT 'gateway.param.read' AS name, 'gateway' AS service, 'param' AS resource, 'read' AS action, 'Read parameters' AS description UNION ALL
  SELECT 'gateway.param.create', 'gateway', 'param', 'create', 'Create parameters' UNION ALL
  SELECT 'gateway.param.update', 'gateway', 'param', 'update', 'Update parameters' UNION ALL
  SELECT 'gateway.param.delete', 'gateway', 'param', 'delete', 'Delete parameters'
)
INSERT INTO public.permission (id, name, service, resource, action, is_active, is_high_risk, description)
SELECT gen_random_uuid(), perms.name, perms.service, perms.resource, perms.action, TRUE, FALSE, perms.description
FROM perms
WHERE NOT EXISTS (SELECT 1 FROM public.permission p WHERE p.name = perms.name);

