package client

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
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

	query := `
		SELECT rm.id, rm.menu_id, m.name AS menu_name, rm.role_id, m.route AS menu_route
		FROM role_menu rm
		JOIN menu m ON rm.menu_id = m.id
		WHERE rm.role_id = ? AND m.is_active = TRUE
		ORDER BY m.sort_order ASC`

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

	args = append(args, req.ID, req.RoleID, req.MenuID, req.CreatedAt)
	query := "INSERT INTO role_menu (id, role_id, menu_id, created_at) VALUES (?, ?, ?, ?)"

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
	searchFields := []string{"m.name", "r.name"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM role_menu rm JOIN menu m ON rm.menu_id = m.id JOIN role r ON rm.role_id = r.id" + whereClause
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
		"id":         "rm.id",
		"menu_id":    "rm.menu_id",
		"menu_name":  "m.name",
		"role_id":    "rm.role_id",
		"role_name":  "r.name",
		"menu_route": "m.route",
		"created_at": "rm.created_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "map.id")

	sb := strings.Builder{}
	sb.WriteString("SELECT rm.id, rm.menu_id, m.name AS menu_name, r.name AS role_name, rm.role_id, m.route AS menu_route, rm.created_at FROM role_menu rm JOIN menu m ON rm.menu_id = m.id JOIN role r ON rm.role_id = r.id")
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
	searchFields := []string{"name", "route"}
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
		"menu_name":  "name",
		"menu_route": "route",
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

	args = append(args, req.ID, req.MenuKey, req.MenuName, req.MenuRoute, req.Icon, req.ParentID, req.SortOrder, req.FeatureKey, req.IsActive, req.CreatedAt, req.UpdatedAt)
	query := "INSERT INTO menu (id, menu_key, name, route, icon, parent_id, sort_order, feature_key, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

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

	// Build base query for counting (without Preload)
	countQuery := r.db.WithContext(ctx).Table("role")

	// Build base query with Preload for Institution
	query := r.db.WithContext(ctx).Table("role").Preload("Institution")

	// Apply search filter to both queries
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		whereClause := "name ILIKE ? OR description ILIKE ?"
		countQuery = countQuery.Where(whereClause, searchPattern, searchPattern)
		query = query.Where(whereClause, searchPattern, searchPattern)
	}

	// Get total count for pagination (using count query without Preload)
	var totalCount int64
	countResult := countQuery.Count(&totalCount)
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
		"id":               "id",
		"role_name":        "name",
		"role_desc":        "description",
		"is_administrator": "is_administrator",
		"is_active":        "is_active",
		"created_at":       "created_at",
		"updated_at":       "updated_at",
		"institution_id":   "institution_id",
	}

	// Apply sorting
	sortBy := filter.SortBy
	sortOrder := filter.SortOrder
	actualField, ok := allowedSortFields[sortBy]
	if !ok {
		actualField = "id"
	}
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}
	query = query.Order(actualField + " " + sortOrder)

	// Apply pagination
	if pagination != nil && pagination.Limit > 0 {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Limit(pagination.Limit).Offset(offset)
	}

	utils.LogEvent(span, "Query", "GORM Query with Preload")

	result := query.Debug().Find(&response)
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

	var response model.Role

	err := r.db.WithContext(ctx).Table("role").Preload("Institution").Where("id = ?", roleID).First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogEventError(span, errors.New("Role not found"))
			return nil, model.ThrowError(http.StatusNotFound, errors.New("Role not found"))
		}
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return &response, nil
}

func (r *RoleClient) CreateNewRole(ctx context.Context, req *model.Role) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewRole")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	// Use InstitutionID if set, otherwise use Institution.ID
	institutionID := req.InstitutionID
	if institutionID == nil && req.Institution != nil && req.Institution.ID != "" {
		id := req.Institution.ID
		institutionID = &id
	}

	var args []interface{}

	args = append(args, req.ID, req.RoleName, req.RoleDesc, req.Scope, req.CreatedAt, req.UpdatedAt, req.IsActive, req.IsAdministrator, institutionID)
	query := "INSERT INTO role (id, name, description, scope, created_at, updated_at, is_active, is_administrator, institution_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

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

	args = append(args, req.RoleName, req.RoleDesc, req.Scope, req.IsAdministrator, req.IsActive, req.UpdatedAt, req.ID)
	query := "UPDATE role SET name = ?, description = ?, scope = ?, is_administrator = ?, is_active = ?, updated_at = ? WHERE id = ?"

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

	args = append(args, req.RoleID, req.MenuID, req.ID)
	query := "UPDATE role_menu SET role_id = ?, menu_id = ? WHERE id = ?"

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

	args = append(args, req.MenuKey, req.MenuName, req.MenuRoute, req.Icon, req.ParentID, req.SortOrder, req.FeatureKey, req.IsActive, req.UpdatedAt, req.ID)
	query := "UPDATE menu SET menu_key = ?, name = ?, route = ?, icon = ?, parent_id = ?, sort_order = ?, feature_key = ?, is_active = ?, updated_at = ? WHERE id = ?"

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

	query := "DELETE FROM role_menu WHERE id = ?"

	err := r.db.Exec(query, id).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Role Mapping")

	return nil
}
