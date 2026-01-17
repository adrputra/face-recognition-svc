package client

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type InterfaceInstitutionClient interface {
	GetAllInstitutions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Institution, *model.Pagination, error)
	GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error)
	CreateNewInstitution(ctx context.Context, institution *model.Institution) error
	UpdateInstitution(ctx context.Context, institution *model.Institution) error
	DeleteInstitution(ctx context.Context, id string) error
}

type InstitutionClient struct {
	db *gorm.DB
}

func NewInstitutionClient(db *gorm.DB) InterfaceInstitutionClient {
	return &InstitutionClient{db: db}
}

func (c *InstitutionClient) GetAllInstitutions(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Institution, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllInstitutions")
	defer span.Finish()

	// Build WHERE clause for search
	searchFields := []string{"name", "code", "address", "email", "phone_number"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM institution" + whereClause
	countResult := c.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount)
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

	var response []*model.Institution

	// Define allowed sort fields
	allowedSortFields := map[string]string{
		"id":           "id",
		"name":         "name",
		"code":         "code",
		"address":      "address",
		"email":        "email",
		"phone_number": "phone_number",
		"created_at":   "created_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "id")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM institution")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}
	query := sb.String()

	utils.LogEvent(span, "Query", query)

	result := c.db.Debug().WithContext(ctx).Raw(query).Scan(&response)
	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	if response == nil {
		return nil, nil, model.ThrowError(http.StatusInternalServerError, errors.New("Data Not Found"))
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
}

func (c *InstitutionClient) GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetInstitutionByID")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)
	var response *model.Institution

	query := "SELECT * FROM institution WHERE id = ?"
	utils.LogEvent(span, "Query", query)

	err := c.db.Debug().WithContext(ctx).Raw(query, id).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)
	return response, nil
}

func (c *InstitutionClient) CreateNewInstitution(ctx context.Context, institution *model.Institution) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", institution)

	var args []interface{}

	args = append(args, institution.ID, institution.Name, institution.Code, institution.Address, institution.PhoneNumber, institution.Email, institution.IsActive, institution.CreatedAt, institution.UpdatedAt)
	err := c.db.Debug().WithContext(ctx).Exec("INSERT INTO institution (id, name, code, address, phone_number, email, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Institution")
	return nil
}

func (c *InstitutionClient) UpdateInstitution(ctx context.Context, institution *model.Institution) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", institution)

	var args []interface{}
	args = append(args, institution.Name, institution.Code, institution.Address, institution.PhoneNumber, institution.Email, institution.IsActive, institution.UpdatedAt, institution.ID)
	err := c.db.Debug().WithContext(ctx).Exec("UPDATE institution SET name = ?, code = ?, address = ?, phone_number = ?, email = ?, is_active = ?, updated_at = ? WHERE id = ?", args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Institution")
	return nil
}

func (c *InstitutionClient) DeleteInstitution(ctx context.Context, id string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)

	err := c.db.Debug().WithContext(ctx).Exec("DELETE FROM institution WHERE id = ?", id).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Institution")
	return nil
}
