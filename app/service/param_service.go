package service

import (
	"errors"
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceParamService interface {
	GetParameterByKey(e echo.Context) error
	GetAllParam(e echo.Context) error
	InsertNewParam(e echo.Context) error
	UpdateParam(e echo.Context) error
	DeleteParam(e echo.Context) error
}

type ParamService struct {
	uc controller.InterfaceParamController
}

func NewParamService(uc controller.InterfaceParamController) *ParamService {
	return &ParamService{uc: uc}
}

func (s *ParamService) GetParameterByKey(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetParameterByKey")
	defer span.Finish()

	key := e.Param("id")
	if key == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	utils.LogEvent(span, "Request", key)

	res, err := s.uc.GetParameterByKey(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Parameter By Key",
		Data:    res,
	})
}

func (s *ParamService) GetAllParam(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllParam")
	defer span.Finish()

	res, err := s.uc.GetAllParam(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All Param",
		Data:    res,
	})
}

func (s *ParamService) InsertNewParam(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "InsertNewParam")
	defer span.Finish()

	var param *model.Param
	if err := e.Bind(&param); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	err := s.uc.InsertNewParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Insert New Param",
		Data:    param,
	})
}

func (s *ParamService) UpdateParam(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateParam")
	defer span.Finish()

	var param *model.Param
	if err := e.Bind(&param); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	err := s.uc.UpdateParam(ctx, param)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Param",
		Data:    param,
	})
}

func (s *ParamService) DeleteParam(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteParam")
	defer span.Finish()

	key := e.Param("id")

	utils.LogEvent(span, "Request", key)

	err := s.uc.DeleteParam(ctx, key)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete Param",
		Data:    nil,
	})
}
