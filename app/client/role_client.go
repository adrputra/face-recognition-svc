package client

import (
	"context"
	"errors"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type InterfaceRoleClient interface {
	GetMenuRoleMapping(ctx context.Context, roleID string) ([]*model.MenuRoleMapping, error)
	CreateNewRoleMapping(ctx context.Context, role *model.MenuRoleMapping) error
	GetAllRoleMapping(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.MenuRoleMapping, *model.Pagination, error)
	UpdateRoleMapping(ctx context.Context, req *model.MenuRoleMapping) error
	DeleteRoleMapping(ctx context.Context, id string) error

	GetAllMenu(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Menu, *model.Pagination, error)
	CreateNewMenu(ctx context.Context, request *model.Menu) error
	UpdateMenu(ctx context.Context, request *model.Menu) error
	DeleteMenu(ctx context.Context, menuID string) error

	GetAllRole(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Role, *model.Pagination, error)
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

	query := "SELECT map.id, map.menu_id, menu.menu_name, map.role_id, menu.menu_route, map.access_method FROM menu_mapping AS map LEFT JOIN menu ON map.menu_id = menu.id LEFT JOIN role ON map.role_id = role.id WHERE role_id = ? ORDER BY map.id ASC"

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

func (r *RoleClient) GetAllRoleMapping(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.MenuRoleMapping, *model.Pagination, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetAllRoleMapping")
	defer span.Finish()

	// Build WHERE clause for search
	searchFields := []string{"menu.menu_name", "role.role_name"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM menu_mapping AS map JOIN menu ON map.menu_id = menu.id JOIN role ON map.role_id = role.id" + whereClause
	countResult := r.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount)
	if countResult.Error != nil {
		utils.LogEventError(span, countResult.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, countResult.Error)
	}

	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	var response []*model.MenuRoleMapping

	// Define allowed sort fields with their actual column names
	allowedSortFields := map[string]string{
		"id":            "map.id",
		"menu_id":       "map.menu_id",
		"menu_name":     "menu.menu_name",
		"role_id":       "map.role_id",
		"role_name":     "role.role_name",
		"menu_route":    "menu.menu_route",
		"access_method": "map.access_method",
		"created_at":    "map.created_at",
		"updated_at":    "map.updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "map.id")

	sb := strings.Builder{}
	sb.WriteString("SELECT map.id, map.menu_id, menu.menu_name, role.role_name, map.role_id, menu.menu_route, map.access_method, map.created_at, map.updated_at, map.created_by, map.updated_by FROM menu_mapping AS map JOIN menu ON map.menu_id = menu.id JOIN role ON map.role_id = role.id")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}
	query := sb.String()

	utils.LogEvent(span, "Query", query)

	result := r.db.Debug().WithContext(ctx).Raw(query).Scan(&response)
	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
}

func (r *RoleClient) GetAllMenu(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Menu, *model.Pagination, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetAllMenu")
	defer span.Finish()

	// Build WHERE clause for search
	searchFields := []string{"menu_name", "menu_route"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM menu" + whereClause
	countResult := r.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount)
	if countResult.Error != nil {
		utils.LogEventError(span, countResult.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, countResult.Error)
	}

	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	var response []*model.Menu

	// Define allowed sort fields
	allowedSortFields := map[string]string{
		"id":         "id",
		"menu_name":  "menu_name",
		"menu_route": "menu_route",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "id")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM menu")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}
	query := sb.String()

	utils.LogEvent(span, "Query", query)

	result := r.db.Debug().WithContext(ctx).Raw(query).Scan(&response)
	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
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

func (r *RoleClient) GetAllRole(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Role, *model.Pagination, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetAllRole")
	defer span.Finish()

	// Build WHERE clause for search
	searchFields := []string{"role_name", "role_desc"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM role" + whereClause
	countResult := r.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount)
	if countResult.Error != nil {
		utils.LogEventError(span, countResult.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, countResult.Error)
	}

	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	var response []*model.Role

	// Define allowed sort fields
	allowedSortFields := map[string]string{
		"id":         "id",
		"role_name":  "role_name",
		"role_desc":  "role_desc",
		"level":      "level",
		"is_active":  "is_active",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "id")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM role")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}
	query := sb.String()

	utils.LogEvent(span, "Query", query)

	result := r.db.Debug().WithContext(ctx).Raw(query).Scan(&response)
	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
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
