package client

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InterfacePermissionClient interface {
	GetAllPermissions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Permission, *model.Pagination, error)
	CreatePermission(ctx context.Context, permission *model.Permission) error
	UpdatePermission(ctx context.Context, request *model.PermissionUpdateRequest) error
	AssignRolePermissions(ctx context.Context, request *model.RolePermissionAssignment) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error)
}

type PermissionClient struct {
	db *gorm.DB
}

func NewPermissionClient(db *gorm.DB) *PermissionClient {
	return &PermissionClient{db: db}
}

func (c *PermissionClient) GetAllPermissions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Permission, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllPermissions")
	defer span.Finish()

	searchFields := []string{"name", "service", "resource", "action"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM permission" + whereClause
	if err := c.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	allowedSortFields := map[string]string{
		"name":       "name",
		"service":    "service",
		"resource":   "resource",
		"action":     "action",
		"is_active":  "is_active",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "name")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM permission")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}

	var response []*model.Permission
	if err := c.db.Debug().WithContext(ctx).Raw(sb.String()).Scan(&response).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	return response, pagination, nil
}

func (c *PermissionClient) CreatePermission(ctx context.Context, permission *model.Permission) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreatePermission")
	defer span.Finish()

	utils.LogEvent(span, "Request", permission)

	query := `
		INSERT INTO permission (id, name, service, resource, action, is_active, is_high_risk, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	args := []interface{}{
		permission.ID,
		permission.Name,
		permission.Service,
		permission.Resource,
		permission.Action,
		permission.IsActive,
		permission.IsHighRisk,
		permission.Description,
		permission.CreatedAt,
		permission.UpdatedAt,
	}

	if err := c.db.Exec(query, args...).Error; err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c *PermissionClient) UpdatePermission(ctx context.Context, request *model.PermissionUpdateRequest) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdatePermission")
	defer span.Finish()

	if request.ID == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("id is required"))
	}

	setClauses := []string{}
	args := []interface{}{}
	if request.IsActive != nil {
		setClauses = append(setClauses, "is_active = ?")
		args = append(args, *request.IsActive)
	}
	if request.IsHighRisk != nil {
		setClauses = append(setClauses, "is_high_risk = ?")
		args = append(args, *request.IsHighRisk)
	}
	if request.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *request.Description)
	}
	if len(setClauses) == 0 {
		return model.ThrowError(http.StatusBadRequest, errors.New("no updatable fields provided"))
	}

	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, request.ID)

	query := "UPDATE permission SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	if err := c.db.Exec(query, args...).Error; err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}

	return nil
}

func (c *PermissionClient) AssignRolePermissions(ctx context.Context, request *model.RolePermissionAssignment) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: AssignRolePermissions")
	defer span.Finish()

	if request.RoleID == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("role_id is required"))
	}
	if len(request.PermissionIDs) == 0 {
		return model.ThrowError(http.StatusBadRequest, errors.New("permission_ids is required"))
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM role_permission WHERE role_id = ?", request.RoleID).Error; err != nil {
			return err
		}
		for _, permissionID := range request.PermissionIDs {
			query := "INSERT INTO role_permission (id, role_id, permission_id, created_at) VALUES (gen_random_uuid(), ?, ?, ?)"
			if err := tx.Exec(query, request.RoleID, permissionID, time.Now()).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *PermissionClient) GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetRolePermissions")
	defer span.Finish()

	var response []*model.Permission
	query := `
		SELECT p.*
		FROM permission p
		JOIN role_permission rp ON rp.permission_id = p.id
		WHERE rp.role_id = ?
		ORDER BY p.name ASC`
	if err := c.db.Debug().WithContext(ctx).Raw(query, roleID).Scan(&response).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, model.ThrowError(http.StatusInternalServerError, err)
	}
	return response, nil
}
