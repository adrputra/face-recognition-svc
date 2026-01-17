package client

import (
	"context"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type InterfaceParamClient interface {
	GetParameterByKey(ctx context.Context, key string) (*model.Param, error)
	GetAllParam(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Param, *model.Pagination, error)
	InsertNewParam(ctx context.Context, param *model.Param) error
	UpdateParam(ctx context.Context, param *model.Param) error
	DeleteParam(ctx context.Context, key string) error
}

type ParamClient struct {
	db *gorm.DB
}

func NewParamClient(db *gorm.DB) *ParamClient {
	return &ParamClient{db: db}
}

func (c *ParamClient) GetParameterByKey(ctx context.Context, key string) (*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetParameterByKey")
	defer span.Finish()

	var result *model.Param

	query := "SELECT * FROM parameter WHERE id = ?"
	err := c.db.Debug().WithContext(ctx).Raw(query, key).Scan(&result).Error

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", result)

	return result, nil
}

func (c *ParamClient) GetAllParam(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Param, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllParam")
	defer span.Finish()

	// Build WHERE clause for search
	searchFields := []string{"id", "value", "description"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM parameter" + whereClause
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

	var result []*model.Param

	// Define allowed sort fields
	allowedSortFields := map[string]string{
		"id":          "id",
		"value":       "value",
		"description": "description",
		"updated_at":  "updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "id")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM parameter")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}
	query := sb.String()

	utils.LogEvent(span, "Query", query)

	resultQuery := c.db.Debug().WithContext(ctx).Raw(query).Scan(&result)
	if resultQuery.Error != nil {
		utils.LogEventError(span, resultQuery.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, resultQuery.Error)
	}

	utils.LogEvent(span, "Response", result)
	utils.LogEvent(span, "Pagination", pagination)

	return result, pagination, nil
}

func (c *ParamClient) InsertNewParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: InsertNewParam")
	defer span.Finish()

	var args []interface{}

	args = append(args, param.Key, param.Value, param.Description, param.UpdatedAt, param.UpdatedBy)
	query := "INSERT INTO parameter (id, value, description, updated_at, updated_by) VALUES (?, ?, ?, ?, ?)"
	result := c.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	utils.LogEvent(span, "Response", "Success Insert New Param")

	return nil
}

func (c *ParamClient) UpdateParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateParam")
	defer span.Finish()

	var args []interface{}

	args = append(args, param.Value, param.UpdatedAt, param.UpdatedBy, param.Key)
	query := "UPDATE parameter SET value = ?, description = ?, updated_at = ?, updated_by = ? WHERE id = ?"
	result := c.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	utils.LogEvent(span, "Response", "Success Update Param")

	return nil
}

func (c *ParamClient) DeleteParam(ctx context.Context, key string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteParam")
	defer span.Finish()

	query := "DELETE FROM parameter WHERE id = ?"

	err := c.db.Debug().WithContext(ctx).Exec(query, key).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Param")

	return nil
}
