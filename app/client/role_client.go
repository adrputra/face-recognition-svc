package client

import (
	"context"
	"errors"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type InterfaceRoleClient interface {
	GetMenuRoleMapping(ctx context.Context, roleID string) ([]*model.MenuRoleMapping, error)
	CreateNewRoleMapping(ctx context.Context, role *model.MenuRoleMapping) error
	GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error)
	UpdateRoleMapping(ctx context.Context, req *model.MenuRoleMapping) error
	DeleteRoleMapping(ctx context.Context, id string) error

	GetAllMenu(ctx context.Context) ([]*model.Menu, error)
	CreateNewMenu(ctx context.Context, request *model.Menu) error
	UpdateMenu(ctx context.Context, request *model.Menu) error
	DeleteMenu(ctx context.Context, menuID string) error

	GetAllRole(ctx context.Context) ([]*model.Role, error)
	GetRoleByID(ctx context.Context, roleID string) (*model.Role, error)
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

	args = append(args, req.Id, req.RoleID, req.MenuID, req.AccessMethod, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy)
	query := "INSERT INTO menu_mapping (id, role_id, menu_id, access_method, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	err := r.db.Exec(query, args...).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Duplicate entry
				utils.LogEventError(span, errors.New("Role Mapping Already Exists"))
				return model.ThrowError(http.StatusBadRequest, errors.New("Role Mapping Already Exists"))
			}
		}
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Role")

	return nil
}

func (r *RoleClient) GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetAllRoleMapping")
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

func (r *RoleClient) GetRoleByID(ctx context.Context, roleID string) (*model.Role, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetRoleByID")
	defer span.Finish()

	utils.LogEvent(span, "Request", roleID)

	var response *model.Role

	query := "SELECT * FROM role WHERE id = ?"

	err := r.db.Debug().Raw(query, roleID).Scan(&response).Error
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

	args = append(args, req.Id, req.RoleName, req.RoleDesc, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy, req.IsActive, req.Level)
	query := "INSERT INTO role (id, role_name, role_desc, created_at, updated_at, created_by, updated_by, is_active, level) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

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

	args = append(args, req.RoleName, req.Level, req.UpdatedAt, req.UpdatedBy, req.Id)
	query := "UPDATE role SET role_name = ?, level = ?, updated_at = ?, updated_by = ? WHERE id = ?"

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

	err := r.db.Debug().Exec(query, args...)
	if err.Error != nil {
		utils.LogEventError(span, err.Error)
		return err.Error
	}

	if err.RowsAffected == 0 {
		utils.LogEventError(span, errors.New("Data Not Found"))
		return errors.New("Data Not Found")
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

func (r *RoleClient) DeleteRoleMapping(ctx context.Context, id string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteRoleMapping")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)

	query := "DELETE FROM menu_mapping WHERE id = ?"

	err := r.db.Exec(query, id).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Role Mapping")

	return nil
}
