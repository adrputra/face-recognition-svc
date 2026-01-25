package client

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InterfaceFeatureClient interface {
	GetAllFeatures(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Feature, *model.Pagination, error)
	CreateFeature(ctx context.Context, feature *model.Feature) error
	UpdateFeature(ctx context.Context, feature *model.Feature) error
	SetInstitutionFeature(ctx context.Context, request *model.InstitutionFeatureRequest) error
	GetInstitutionFeatures(ctx context.Context, institutionID string) ([]*model.InstitutionFeature, error)
}

type FeatureClient struct {
	db *gorm.DB
}

func NewFeatureClient(db *gorm.DB) *FeatureClient {
	return &FeatureClient{db: db}
}

func (c *FeatureClient) GetAllFeatures(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.Feature, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllFeatures")
	defer span.Finish()

	searchFields := []string{"feature_key", "name", "feature_type"}
	whereClause := utils.BuildSearchWhereClause(filter.Search, searchFields)

	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM feature" + whereClause
	if err := c.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	allowedSortFields := map[string]string{
		"feature_key":     "feature_key",
		"name":            "name",
		"feature_type":    "feature_type",
		"default_enabled": "default_enabled",
		"created_at":      "created_at",
		"updated_at":      "updated_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "feature_key")

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM feature")
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}

	var response []*model.Feature
	if err := c.db.Debug().WithContext(ctx).Raw(sb.String()).Scan(&response).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	return response, pagination, nil
}

func (c *FeatureClient) CreateFeature(ctx context.Context, feature *model.Feature) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateFeature")
	defer span.Finish()

	query := `
		INSERT INTO feature (id, feature_key, name, description, feature_type, default_enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	args := []interface{}{
		feature.ID,
		feature.FeatureKey,
		feature.Name,
		feature.Description,
		feature.FeatureType,
		feature.DefaultEnabled,
		feature.CreatedAt,
		feature.UpdatedAt,
	}

	if err := c.db.Exec(query, args...).Error; err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c *FeatureClient) UpdateFeature(ctx context.Context, feature *model.Feature) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateFeature")
	defer span.Finish()

	if feature.ID == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("id is required"))
	}

	query := `
		UPDATE feature
		SET name = ?, description = ?, feature_type = ?, default_enabled = ?, updated_at = ?
		WHERE id = ?`
	args := []interface{}{
		feature.Name,
		feature.Description,
		feature.FeatureType,
		feature.DefaultEnabled,
		feature.UpdatedAt,
		feature.ID,
	}

	if err := c.db.Exec(query, args...).Error; err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c *FeatureClient) SetInstitutionFeature(ctx context.Context, request *model.InstitutionFeatureRequest) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: SetInstitutionFeature")
	defer span.Finish()

	if request.InstitutionID == "" || request.FeatureKey == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("institution_id and feature_key are required"))
	}

	query := `
		INSERT INTO institution_feature (id, institution_id, feature_key, is_enabled, created_at, updated_at)
		VALUES (gen_random_uuid(), ?, ?, ?, ?, ?)
		ON CONFLICT (institution_id, feature_key) DO UPDATE
		SET is_enabled = EXCLUDED.is_enabled, updated_at = EXCLUDED.updated_at`
	args := []interface{}{
		request.InstitutionID,
		request.FeatureKey,
		request.IsEnabled,
		time.Now(),
		time.Now(),
	}

	if err := c.db.Exec(query, args...).Error; err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c *FeatureClient) GetInstitutionFeatures(ctx context.Context, institutionID string) ([]*model.InstitutionFeature, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetInstitutionFeatures")
	defer span.Finish()

	var response []*model.InstitutionFeature
	query := "SELECT * FROM institution_feature WHERE institution_id = ?"
	if err := c.db.Debug().WithContext(ctx).Raw(query, institutionID).Scan(&response).Error; err != nil {
		utils.LogEventError(span, err)
		return nil, model.ThrowError(http.StatusInternalServerError, err)
	}
	return response, nil
}
