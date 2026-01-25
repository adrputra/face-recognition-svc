---
title: Gateway RBAC Technical PRD
version: 1.0
service: gateway
last_updated: 2026-01-17
---

# Gateway RBAC Technical PRD

This document defines the RBAC design and behavior implemented by the Gateway Service, derived from the current database schema and service architecture. It is intended for backend engineers and platform owners who need a precise, technical description of authorization, tenancy, and UI navigation resolution.

## 1) Goals and Non-Goals

### Goals
- Single authorization authority at the Gateway boundary.
- Users can belong to multiple institutions and hold multiple roles per institution.
- Permission checks are deterministic, cache-friendly, and enforced at API boundaries.
- UI menu resolution is derived from roles and feature toggles, never used for authorization.
- Extendable to new services without schema rewrites.

### Non-Goals
- No direct permissions on users.
- No authorization based on menu visibility.
- No implicit cross-institution access.

## 2) Core Concepts

### Tenancy Boundary
The institution is the strict tenant boundary. Every permission check must be scoped to the active institution.

### Role Assignment
Roles are assigned **per institution** using a membership link. Users may have:
- Multiple roles per institution.
- Membership in multiple institutions.

### Permissions
Permissions follow an immutable name format:
```
<service>.<resource>.<action>
```
Permissions are owned by roles and are enforced at API/gRPC boundaries.

### Menus
Menus represent navigation only. Visibility is determined by:
```
user_roles → role_menus → menus → features.enabled
```
Menus never grant permissions.

## 3) Database Schema (RBAC Scope)

### Identity & Tenancy
- `user`
  - `id`, `username`, `email`, `password_hash`, `full_name`, `short_name`, `is_active`
  - `profile_photo`, `cover_photo`
- `institution`
  - `id`, `name`, `code`, `address`, `phone_number`, `email`, `is_active`
- `user_institution`
  - `user_id`, `institution_id`, `status`, `joined_at`, `left_at`

### RBAC Core
- `role`
  - `id`, `name`, `description`, `scope` (`system` | `institution`)
  - `institution_id` (NULL for `system` scope)
  - `is_active`, `is_administrator`
- `user_role`
  - `user_id`, `institution_id`, `role_id`
  - Foreign key to `user_institution` ensures membership exists
- `permission`
  - `name` (unique), `service`, `resource`, `action`
  - `is_active`, `is_high_risk`
  - Immutable identity enforced by trigger
- `role_permission`
  - `role_id`, `permission_id`

### UI Layer
- `menu`
  - `menu_key`, `name`, `route`, `icon`, `parent_id`, `sort_order`
  - `feature_key`, `is_active`
- `role_menu`
  - `role_id`, `menu_id`

### Feature Toggles
- `feature`
  - `feature_key`, `feature_type`, `default_enabled`
- `institution_feature`
  - `institution_id`, `feature_key`, `is_enabled`

### Audit
- `audit_log`
  - `actor_user_id`, `institution_id`, `permission_name`, `action`, `entity_type`, `entity_id`
  - `request_id`, `ip_address`, `user_agent`, `metadata`

## 4) Authorization Flow

### 4.1 Authentication
1. User logs in with `username`, `password`, and `institution_id`.
2. Gateway validates credentials and membership.
3. Gateway issues JWT with:
   - `user_id`
   - `username`
   - `role_ids` (role list for current institution)
   - `institution_id`

### 4.2 Authorization (REST)
1. Request hits Gateway with JWT.
2. Middleware extracts JWT claims.
3. If `app-permission` header is present, middleware checks:
   - `permission` exists and active
   - user has `role_permission` via any role in `role_ids`
4. If permission is missing → `403`.

**Note:** No menu-based logic is used for authorization. The request is authorized purely by permission name.

### 4.3 Authorization (gRPC)
Gateway injects metadata headers:
- `user_id`
- `username`
- `role_ids` (comma-separated)
- `institution_id`

Downstream gRPC interceptors enforce permission by method mapping or use Gateway pre-validated permission context.

## 5) Role Resolution Algorithm

Given `(user_id, institution_id)`:
```
1. Verify membership: user_institution exists and status = active.
2. Resolve roles:
   roles = SELECT role_id FROM user_role
           WHERE user_id = ? AND institution_id = ?
3. Resolve permissions:
   perms = SELECT p.name FROM permission p
           JOIN role_permission rp ON rp.permission_id = p.id
           WHERE rp.role_id IN roles AND p.is_active = TRUE
4. Use perms for authorization decisions.
```

## 6) Menu Rendering Algorithm

Given `(user_id, institution_id)`:
```
1. roles = user_role(user_id, institution_id)
2. menus = role_menu(roles) JOIN menu WHERE menu.is_active = TRUE
3. features = institution_feature(institution_id)
4. filter menus:
   - if menu.feature_key is NULL → keep
   - else keep only if feature enabled
5. build hierarchy by parent_id + sort_order
```

## 7) Feature Toggles

### Effective Feature Evaluation
```
if institution_feature exists:
    enabled = institution_feature.is_enabled
else:
    enabled = feature.default_enabled
```

### Behavior
- `feature_type = menu` controls visibility
- `feature_type = permission` can gate permission creation or global enable/disable
- `feature_type = system` reserved for platform behavior

## 8) Audit and Governance

### Audit Events
Log on:
- Authentication attempts
- Role assignment updates
- Permission assignment updates
- Permission enforcement failures
- High-risk permission usage

### Immutable Controls
`permission` identity fields cannot be changed due to trigger enforcement.

## 9) Cache Strategy (Recommended)

### Cache Keys
- `permset:{user_id}:{institution_id}` → set of permissions
- `roleset:{user_id}:{institution_id}` → role list

### Invalidations
Invalidate on:
- role assignment changes
- role_permission changes
- permission activation changes

## 10) End-to-End Example

**Scenario:** User `U1` logs into institution `I1`.

1. Membership exists in `user_institution`.
2. Roles resolved:
   - `role: teacher`
   - `role: class_manager`
3. Permissions resolved:
   - `class.grade.create`
   - `presence.attendance.mark`
4. Request requires `presence.attendance.mark` → allowed.
5. Menu resolution returns tree filtered by `role_menu` and `feature` flags.

## 11) API Boundary Requirements

For every protected REST route:
- The route must declare required permission name in `app-permission` header.
- Gateway must enforce permission before hitting downstream services.
- JWT must include `institution_id`.

For every gRPC call:
- Gateway sets metadata (see 4.3).
- Downstream service enforces permission for method or relies on Gateway validation.

## 12) Data Integrity Constraints (Highlights)

- `user_role` requires `user_institution` link.
- `role` enforces scope vs `institution_id` constraint.
- `permission` name pattern enforced by regex.
- All tables use UUID primary keys.

## 13) Migration Notes

- Initial RBAC schema is created in `000013_rbac_redesign.up.sql`.
- Patch migration `000014_sync_rbac_columns.up.sql` aligns missing columns for compatibility.

## 14) Operational Checklist

- Ensure database migrations are applied in order.
- Seed at least one `institution` and `system` role.
- Seed required permissions and assign to roles.
- Ensure JWT secret is configured and consistent.
- Ensure `app-permission` header is set by API clients or gateway routing layer.

## 15) Feature Usage Guide

### Definition
Feature is a **toggleable capability** used to control UI navigation (and optionally other system behavior). It is **not** a permission and does **not** grant API access.

### Where it applies
- `menu.feature_key` links a menu item to a feature toggle.
- `institution_feature` overrides feature state per institution.
- If no override exists, `feature.default_enabled` is used.

### Recommended Naming
Use a namespaced, purpose-specific key. Examples:
- `menu.dashboard`
- `menu.user_management`
- `menu.attendance`
- `system.beta_reports`

### Example: Enable menu per institution
1) Create feature:
```
feature_key: "menu.dashboard"
feature_type: "menu"
default_enabled: true
```

2) Create or update menu:
```
menu_key: "dashboard"
feature_key: "menu.dashboard"
```

3) Override per institution:
```
institution_id: "12345"
feature_key: "menu.dashboard"
is_enabled: true
```

### Rules
- Menus with `feature_key = NULL` are always eligible.
- Menu visibility is derived from `role_menu` **and** feature state.
- Feature toggles never replace permission checks.
