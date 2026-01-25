package service

import (
	"errors"
	"face-recognition-svc/gateway/app/controller"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceFeatureService interface {
	GetAllFeatures(e echo.Context) error
	CreateFeature(e echo.Context) error
	UpdateFeature(e echo.Context) error
	SetInstitutionFeature(e echo.Context) error
	GetInstitutionFeatures(e echo.Context) error
}

type FeatureService struct {
	fc controller.InterfaceFeatureController
}

func NewFeatureService(fc controller.InterfaceFeatureController) InterfaceFeatureService {
	return &FeatureService{fc: fc}
}

func (s *FeatureService) GetAllFeatures(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllFeatures")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	res, pagination, err := s.fc.GetAllFeatures(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Feature",
		Data:       res,
		Pagination: pagination,
	})
}

func (s *FeatureService) CreateFeature(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateFeature")
	defer span.Finish()

	var request *model.Feature
	if err := e.Bind(&request); err != nil {
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

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create Feature",
		Data:    nil,
	})
}

func (s *FeatureService) UpdateFeature(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateFeature")
	defer span.Finish()

	var request *model.Feature
	if err := e.Bind(&request); err != nil {
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

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Feature",
		Data:    nil,
	})
}

func (s *FeatureService) SetInstitutionFeature(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "SetInstitutionFeature")
	defer span.Finish()

	var request *model.InstitutionFeatureRequest
	if err := e.Bind(&request); err != nil {
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

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Set Institution Feature",
		Data:    nil,
	})
}

func (s *FeatureService) GetInstitutionFeatures(e echo.Context) error {
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

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Institution Features",
		Data:    res,
	})
}
