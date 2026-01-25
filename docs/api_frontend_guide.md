---
title: Gateway API Frontend Guide
version: 1.0
service: gateway
last_updated: 2026-01-17
---

# Gateway API Frontend Guide (AI-Friendly)

This document describes how a frontend client should use the Gateway Service APIs and which **forms** and **data fields** are required to implement the UI. It is written to be both human- and AI-consumable.

## 1) Base and Auth

### Base URL
```
{{endpoint}}
```

### Authentication
- Use **JWT Bearer token** from `POST /api/login`.
- All `/api/service/*` endpoints require `Authorization: Bearer {{token}}`.

### Authorization
- Some endpoints may require `app-permission` header (e.g., protected internal flows).
- Menu visibility is NOT authorization.

## 2) Common Response Shape

```
{
  "code": 200,
  "message": "Success ...",
  "data": <object | array | null>,
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## 3) Forms and Endpoints

### 3.1 Authentication and Session

#### Login (Step 1: credentials only)
**Endpoint**
```
POST /api/login
```
**Form Fields**
- `username` (string, required)
- `password` (string, required)
- `institution_id` (string, optional)

**Response Data**
- If `institution_id` is **not** provided:
  - `institutions` (array of institutions the user belongs to)
  - `token` is empty
- If `institution_id` is provided:
  - `token` (string)
  - `role_ids` (array of role IDs)
  - `institution_id`, `institution_name`
  - `menu_mapping` (array of menu items for UI)

#### Register (Create User)
**Endpoint**
```
POST /api/register
```
**Form Fields**
- `username` (string, required)
- `email` (string, required)
- `password` (string, required)
- `fullname` (string, required)
- `shortname` (string, optional)
- `institution_id` (string, required)
- `role_ids` (array of role IDs, required)

### 3.2 User Management

#### List Users
```
GET /api/service/user?page=1&limit=10&search=...&sort_by=...&sort_order=...
```

#### User Detail
```
GET /api/service/user/detail/:username
```

#### Update User
```
PUT /api/service/user
```
**Form Fields**
- `username` (string, required)
- `email` (string, optional)
- `fullname` (string, optional)
- `shortname` (string, optional)
- `is_active` (bool, optional)
- `institution_id` (string, required to update role assignments)
- `role_ids` (array of role IDs, required if updating roles)

#### Delete User
```
DELETE /api/service/user/:username
```

#### Institution List (for user creation)
```
GET /api/service/user/institutions
```

#### Upload Profile Photo
```
POST /api/service/user/profile-photo
```
**Form Data**
- `file` (file, required)

#### Upload Cover Photo
```
POST /api/service/user/cover-photo
```
**Form Data**
- `file` (file, required)

### 3.3 Institution Management

#### List Institutions
```
GET /api/service/institution
```

#### Institution Detail
```
GET /api/service/institution/:id
```

#### Create Institution
```
POST /api/service/institution
```
**Form Fields**
- `name` (string, required)
- `code` (string, required)
- `address` (string, optional)
- `phone_number` (string, optional)
- `email` (string, optional)
- `is_active` (bool, optional)

#### Update Institution
```
PUT /api/service/institution
```
**Form Fields**
- `id` (string, required)
- `name`, `code`, `address`, `phone_number`, `email`, `is_active`

#### Delete Institution
```
DELETE /api/service/institution/:id
```

### 3.4 Role Management

#### List Roles
```
GET /api/service/role
```

#### Create Role
```
POST /api/service/role/create
```
**Form Fields**
- `role_name` (string, required)
- `role_desc` (string, optional)
- `scope` (string: `system` or `institution`, required)
- `institution_id` (string, required if scope = institution)
- `is_active` (bool, optional)
- `is_administrator` (bool, optional)

### 3.5 Menu Management

#### List Menus
```
GET /api/service/role/menu
```

#### Create Menu
```
POST /api/service/role/menu/create
```
**Form Fields**
- `menu_key` (string, required)
- `menu_name` (string, required)
- `menu_route` (string, optional)
- `icon` (string, optional)
- `parent_id` (string or null)
- `sort_order` (int, optional)
- `feature_key` (string or null)
- `is_active` (bool, optional)

#### Update Menu
```
PUT /api/service/role/menu
```
**Form Fields**
- `id` (string, required)
- `menu_key`, `menu_name`, `menu_route`, `icon`, `parent_id`, `sort_order`, `feature_key`, `is_active`

#### Delete Menu
```
DELETE /api/service/role/menu/:id
```

### 3.6 Role ↔ Menu Mapping

#### List Role Menu Mapping
```
GET /api/service/role/mapping
```

#### Create Role Menu Mapping
```
POST /api/service/role/mapping/create
```
**Form Fields**
- `role_id` (string, required)
- `menu_id` (string, required)

#### Update Role Menu Mapping
```
PUT /api/service/role/mapping
```
**Form Fields**
- `id` (string, required)
- `role_id` (string, required)
- `menu_id` (string, required)

#### Delete Role Menu Mapping
```
DELETE /api/service/role/mapping/:id
```

### 3.7 Permission Management

#### List Permissions
```
GET /api/service/permission
```

#### Create Permission
```
POST /api/service/permission
```
**Form Fields**
- `name` (string, required; format `<service>.<resource>.<action>`)
- `service` (string, required)
- `resource` (string, required)
- `action` (string, required)
- `is_active` (bool, optional)
- `is_high_risk` (bool, optional)
- `description` (string, optional)

#### Update Permission (limited fields)
```
PUT /api/service/permission
```
**Form Fields**
- `id` (string, required)
- `is_active` (bool, optional)
- `is_high_risk` (bool, optional)
- `description` (string, optional)

#### Assign Permissions to Role
```
POST /api/service/permission/assign
```
**Form Fields**
- `role_id` (string, required)
- `permission_ids` (array of permission IDs, required)

#### Get Role Permissions
```
GET /api/service/permission/role/:id
```

### 3.8 Feature Management

#### List Features
```
GET /api/service/feature
```

#### Create Feature
```
POST /api/service/feature
```
**Form Fields**
- `feature_key` (string, required)
- `name` (string, required)
- `description` (string, optional)
- `feature_type` (string: `menu`, `permission`, `system`)
- `default_enabled` (bool)

#### Update Feature
```
PUT /api/service/feature
```
**Form Fields**
- `id` (string, required)
- `name`, `description`, `feature_type`, `default_enabled`

#### Set Institution Feature Override
```
POST /api/service/feature/institution
```
**Form Fields**
- `institution_id` (string, required)
- `feature_key` (string, required)
- `is_enabled` (bool, required)

#### Get Institution Features
```
GET /api/service/feature/institution/:id
```

### 3.9 Dataset Management

#### List Datasets
```
GET /api/service/dataset
```

#### Upload Dataset
```
POST /api/service/dataset
```
**Form Data**
- `username` (string)
- `file` (file, multi)

#### Delete Dataset
```
DELETE /api/service/dataset/:username
```

#### Train Model
```
POST /api/service/dataset/train-model/:institution_id
```

#### Get Last Training
```
GET /api/service/dataset/last-train-model/:institution_id
```

#### Training History
```
POST /api/service/dataset/model-training-history
```
**Form Fields**
- `institution_id` (string)
- `status` (string)
- `is_used` (string)
- `order_by` (string)
- `sort_type` (string)

#### Get Datasets by Username
```
GET /api/service/dataset/:institution-id/:username
```

### 3.10 Parameter Management

#### Get Parameter
```
GET /api/service/param/:key
```

#### List Parameters
```
GET /api/service/param
```

#### Create Parameter
```
POST /api/service/param
```
**Form Fields**
- `key` (string, required)
- `value` (string, required)
- `description` (string, optional)

#### Update Parameter
```
PUT /api/service/param
```
**Form Fields**
- `key` (string, required)
- `value` (string, required)
- `description` (string, optional)

#### Delete Parameter
```
DELETE /api/service/param/:key
```

## 4) UI Page Checklist (Suggested)

- Login page (username, password, institution selector)
- User management (list, create, edit, delete)
- Institution management (list, create, edit)
- Role management (list, create)
- Menu management (list, create, edit, link feature)
- Role-menu mapping (assign menus to roles)
- Permission management (list, create, assign to role)
- Feature management (list, create, toggle per institution)
- Dataset management (upload, list, train)
- Parameters (list, update)

## 5) Notes for AI UI Generation

- Each form should map directly to a single endpoint payload.
- Use explicit fields exactly as defined above (case-sensitive JSON keys).
- Menu creation should allow optional `feature_key` and `parent_id` (nullable).
- Permission name must match the regex pattern from the DB constraint.
