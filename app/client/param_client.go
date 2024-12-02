package client

import (
	"context"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"

	"gorm.io/gorm"
)

type InterfaceParamClient interface {
	GetParameterByKey(ctx context.Context, key string) (*model.Param, error)
	GetAllParam(ctx context.Context) ([]*model.Param, error)
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

func (c *ParamClient) GetAllParam(ctx context.Context) ([]*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetDatasetList")
	defer span.Finish()

	var result []*model.Param

	query := "SELECT * FROM parameter"
	err := c.db.Debug().WithContext(ctx).Raw(query).Scan(&result).Error

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", result)

	return result, nil
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

	err := c.db.Exec(query, key).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Param")

	return nil
}
