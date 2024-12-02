package controller

import (
	"context"
	"encoding/json"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type InterfaceParamController interface {
	GetParameterByKey(ctx context.Context, key string) (*model.Param, error)
	GetAllParam(ctx context.Context) ([]*model.Param, error)
	InsertNewParam(ctx context.Context, param *model.Param) error
	UpdateParam(ctx context.Context, param *model.Param) error
	DeleteParam(ctx context.Context, key string) error
}

type ParamController struct {
	redis  *redis.Client
	client client.InterfaceParamClient
}

func NewParamController(redis *redis.Client, client client.InterfaceParamClient) *ParamController {
	return &ParamController{
		redis:  redis,
		client: client,
	}
}

func (c *ParamController) GetParameterByKey(ctx context.Context, key string) (*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetParameterByKey")
	defer span.Finish()

	utils.LogEvent(span, "Request", key)

	cache := c.redis.Get(ctx, key).Val() // Get string value from Redis
	if cache != "" {
		utils.LogEvent(span, "Redis", cache)

		// Deserialize the cached value into the expected object
		resCache := &model.Param{}
		if err := json.Unmarshal([]byte(cache), resCache); err != nil {
			utils.LogEventError(span, err)
		} else {
			return resCache, nil // Return the cached value
		}
	}

	res, err := c.client.GetParameterByKey(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	// Serialize the response into JSON
	resJSON, err := json.Marshal(res)
	if err != nil {
		utils.LogEventError(span, err)
		return res, nil // Return the result even if caching fails
	}

	// Set the serialized object in Redis with an expiration (e.g., 5 minutes)
	if err := c.redis.Set(ctx, key, resJSON, 6*time.Hour).Err(); err != nil {
		utils.LogEventError(span, err)
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *ParamController) GetAllParam(ctx context.Context) ([]*model.Param, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllParam")
	defer span.Finish()

	utils.LogEvent(span, "Request", "All")

	res, err := c.client.GetAllParam(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *ParamController) InsertNewParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: InsertNewParam")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	param.UpdatedAt = time.Now()
	param.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", param)

	err = c.client.InsertNewParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Insert New Param")

	return nil
}

func (c *ParamController) UpdateParam(ctx context.Context, param *model.Param) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateParam")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	param.UpdatedAt = time.Now()
	param.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", param)

	err = c.client.UpdateParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	newValueJSON, err := json.Marshal(param)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	if err := c.redis.Set(ctx, param.Key, newValueJSON, 6*time.Hour).Err(); err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Param")

	return nil
}

func (c *ParamController) DeleteParam(ctx context.Context, key string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteParam")
	defer span.Finish()

	utils.LogEvent(span, "Request", key)

	err := c.client.DeleteParam(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	err = c.redis.Del(ctx, key).Err()
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Param")

	return nil
}
