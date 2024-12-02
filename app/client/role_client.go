package client

import (
	"context"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"

	"gorm.io/gorm"
)

type InterfaceRoleClient interface {
	GetMenuRoleMapping(ctx context.Context, roleID string) ([]*model.MenuRoleMapping, error)
	CreateNewRoleMapping(ctx context.Context, role *model.MenuRoleMapping) error
	GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error)
	UpdateRoleMapping(ctx context.Context, req *model.MenuRoleMapping) error

	GetAllMenu(ctx context.Context) ([]*model.Menu, error)
	CreateNewMenu(ctx context.Context, request *model.Menu) error
	UpdateMenu(ctx context.Context, request *model.Menu) error
	DeleteMenu(ctx context.Context, menuID string) error

	GetAllRole(ctx context.Context) ([]*model.Role, error)
	CreateNewRole(ctx context.Context, request *model.Role) error
	UpdateRole(ctx context.Context, request *model.Role) error
}

type RoleClient struct {
	db *gorm.DB
}

func NewRoleClient(db *gorm.DB) *RoleClient {
	return &RoleClient{db: db}
}

func (r *RoleClient) GetMenuRoleMapping(ctx context.Context, roleID string) ([]*model.MenuRoleMapping, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetMenuRoleMapping")
	defer span.Finish()

	utils.LogEvent(span, "Request", roleID)

	var response []*model.MenuRoleMapping

	query := "SELECT map.id, map.menu_id, menu.menu_name, map.role_id, menu.menu_route, map.access_method, map.created_at, map.updated_at, map.created_by, map.updated_by FROM menu_mapping AS map JOIN menu ON map.menu_id = menu.id JOIN role ON map.role_id = role.id WHERE role_id = ? ORDER BY map.id ASC"

	utils.LogEvent(span, "Query", query)

	err := r.db.Debug().Raw(query, roleID).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)
	return response, nil
}

func (r *RoleClient) CreateNewRoleMapping(ctx context.Context, req *model.MenuRoleMapping) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewRoleMapping")

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.RoleID, req.MenuID, req.AccessMethod, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy)
	query := "INSERT INTO menu_mapping (role_id, menu_id, access_method, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Role")

	return nil
}

func (r *RoleClient) GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllRoleMapping")
	defer span.Finish()

	var response []*model.MenuRoleMapping

	query := "SELECT map.id, map.menu_id, menu.menu_name, role.role_name, map.role_id, menu.menu_route, map.access_method, map.created_at, map.updated_at, map.created_by, map.updated_by FROM menu_mapping AS map JOIN menu ON map.menu_id = menu.id JOIN role ON map.role_id = role.id ORDER BY map.id ASC"

	err := r.db.Debug().Raw(query).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (r *RoleClient) GetAllMenu(ctx context.Context) ([]*model.Menu, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllMenu")
	defer span.Finish()

	var response []*model.Menu

	query := "SELECT * FROM menu ORDER BY id ASC"

	err := r.db.Debug().Raw(query).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (r *RoleClient) CreateNewMenu(ctx context.Context, req *model.Menu) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewMenu")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.Id, req.MenuName, req.MenuRoute, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy)
	query := "INSERT INTO menu (id, menu_name, menu_route, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Menu")

	return nil
}

func (r *RoleClient) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllRole")
	defer span.Finish()

	var response []*model.Role

	query := "SELECT * FROM role ORDER BY id ASC"

	err := r.db.Debug().Raw(query).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (r *RoleClient) CreateNewRole(ctx context.Context, req *model.Role) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewRole")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.Id, req.RoleName, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy)
	query := "INSERT INTO role (id, role_name, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Role")

	return nil
}

func (r *RoleClient) UpdateRole(ctx context.Context, req *model.Role) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateRole")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.RoleName, req.UpdatedAt, req.UpdatedBy, req.Id)
	query := "UPDATE role SET role_name = ?, updated_at = ?, updated_by = ? WHERE id = ?"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Role")

	return nil
}

func (r *RoleClient) UpdateRoleMapping(ctx context.Context, req *model.MenuRoleMapping) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateRoleMapping")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.AccessMethod, req.UpdatedAt, req.UpdatedBy, req.Id)
	query := "UPDATE menu_mapping SET access_method = ?, updated_at = ?, updated_by = ? WHERE id = ?"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Role Mapping")

	return nil
}

func (r *RoleClient) UpdateMenu(ctx context.Context, req *model.Menu) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateMenu")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var args []interface{}

	args = append(args, req.MenuName, req.MenuRoute, req.UpdatedAt, req.UpdatedBy, req.Id)
	query := "UPDATE menu SET menu_name = ?, menu_route = ?, updated_at = ?, updated_by = ? WHERE id = ?"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Menu")

	return nil
}

func (r *RoleClient) DeleteMenu(ctx context.Context, id string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteMenu")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)

	query := "DELETE FROM menu WHERE id = ?"

	err := r.db.Exec(query, id).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Menu")

	return nil
}
