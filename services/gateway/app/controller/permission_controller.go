package controller

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/client"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type InterfacePermissionController interface {
	GetAllPermissions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Permission, *model.Pagination, error)
	CreatePermission(ctx context.Context, request *model.Permission) error
	UpdatePermission(ctx context.Context, request *model.PermissionUpdateRequest) error
	AssignRolePermissions(ctx context.Context, request *model.RolePermissionAssignment) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error)
}

type PermissionController struct {
	permissionClient client.InterfacePermissionClient
}

func NewPermissionController(permissionClient client.InterfacePermissionClient) *PermissionController {
	return &PermissionController{permissionClient: permissionClient}
}

func (c *PermissionController) GetAllPermissions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Permission, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllPermissions")
	defer span.Finish()

	return c.permissionClient.GetAllPermissions(ctx, pagination, filter)
}

func (c *PermissionController) CreatePermission(ctx context.Context, request *model.Permission) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreatePermission")
	defer span.Finish()

	if request.Name == "" || request.Service == "" || request.Resource == "" || request.Action == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("name, service, resource, action are required"))
	}

	request.ID = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	return c.permissionClient.CreatePermission(ctx, request)
}

func (c *PermissionController) UpdatePermission(ctx context.Context, request *model.PermissionUpdateRequest) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdatePermission")
	defer span.Finish()

	return c.permissionClient.UpdatePermission(ctx, request)
}

func (c *PermissionController) AssignRolePermissions(ctx context.Context, request *model.RolePermissionAssignment) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: AssignRolePermissions")
	defer span.Finish()

	return c.permissionClient.AssignRolePermissions(ctx, request)
}

func (c *PermissionController) GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetRolePermissions")
	defer span.Finish()

	if roleID == "" {
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("role_id is required"))
	}
	return c.permissionClient.GetRolePermissions(ctx, roleID)
}
