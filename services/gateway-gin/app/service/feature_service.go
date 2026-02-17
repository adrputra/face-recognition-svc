package service

import (
	"errors"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/controller"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterfaceFeatureService interface {
	GetAllFeatures(e *gin.Context) error
	CreateFeature(e *gin.Context) error
	UpdateFeature(e *gin.Context) error
	SetInstitutionFeature(e *gin.Context) error
	GetInstitutionFeatures(e *gin.Context) error
}

type FeatureService struct {
	fc controller.InterfaceFeatureController
}

func NewFeatureService(fc controller.InterfaceFeatureController) InterfaceFeatureService {
	return &FeatureService{fc: fc}
}

func (s *FeatureService) GetAllFeatures(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllFeatures")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	res, pagination, err := s.fc.GetAllFeatures(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Feature",
		Data:       res,
		Pagination: pagination,
	})
}

func (s *FeatureService) CreateFeature(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "CreateFeature")
	defer span.Finish()

	var request *model.Feature
	if err := e.ShouldBind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}
	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.fc.CreateFeature(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create Feature",
		Data:    nil,
	})
}

func (s *FeatureService) UpdateFeature(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateFeature")
	defer span.Finish()

	var request *model.Feature
	if err := e.ShouldBind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}
	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.fc.UpdateFeature(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Feature",
		Data:    nil,
	})
}

func (s *FeatureService) SetInstitutionFeature(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "SetInstitutionFeature")
	defer span.Finish()

	var request *model.InstitutionFeatureRequest
	if err := e.ShouldBind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}
	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.fc.SetInstitutionFeature(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Set Institution Feature",
		Data:    nil,
	})
}

func (s *FeatureService) GetInstitutionFeatures(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetInstitutionFeatures")
	defer span.Finish()

	institutionID := e.Param("id")
	if institutionID == "" {
		return utils.LogError(e, errors.New("institution_id shouldn't be empty"), nil)
	}

	res, err := s.fc.GetInstitutionFeatures(ctx, institutionID)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Institution Features",
		Data:    res,
	})
}
