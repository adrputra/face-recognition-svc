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

type InterfaceFeatureController interface {
	GetAllFeatures(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Feature, *model.Pagination, error)
	CreateFeature(ctx context.Context, request *model.Feature) error
	UpdateFeature(ctx context.Context, request *model.Feature) error
	SetInstitutionFeature(ctx context.Context, request *model.InstitutionFeatureRequest) error
	GetInstitutionFeatures(ctx context.Context, institutionID string) ([]*model.InstitutionFeature, error)
}

type FeatureController struct {
	featureClient client.InterfaceFeatureClient
}

func NewFeatureController(featureClient client.InterfaceFeatureClient) *FeatureController {
	return &FeatureController{featureClient: featureClient}
}

func (c *FeatureController) GetAllFeatures(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Feature, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllFeatures")
	defer span.Finish()

	return c.featureClient.GetAllFeatures(ctx, pagination, filter)
}

func (c *FeatureController) CreateFeature(ctx context.Context, request *model.Feature) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreateFeature")
	defer span.Finish()

	if request.FeatureKey == "" || request.Name == "" || request.FeatureType == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("feature_key, name, feature_type are required"))
	}

	request.ID = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	return c.featureClient.CreateFeature(ctx, request)
}

func (c *FeatureController) UpdateFeature(ctx context.Context, request *model.Feature) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateFeature")
	defer span.Finish()

	if request.ID == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("id is required"))
	}

	request.UpdatedAt = time.Now()
	return c.featureClient.UpdateFeature(ctx, request)
}

func (c *FeatureController) SetInstitutionFeature(ctx context.Context, request *model.InstitutionFeatureRequest) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: SetInstitutionFeature")
	defer span.Finish()

	return c.featureClient.SetInstitutionFeature(ctx, request)
}

func (c *FeatureController) GetInstitutionFeatures(ctx context.Context, institutionID string) ([]*model.InstitutionFeature, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetInstitutionFeatures")
	defer span.Finish()

	if institutionID == "" {
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("institution_id is required"))
	}
	return c.featureClient.GetInstitutionFeatures(ctx, institutionID)
}
