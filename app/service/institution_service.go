package service

import (
	"errors"
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceInstitutionService interface {
	GetAllInstitution(e echo.Context) error
	GetInstitutionByID(e echo.Context) error
	CreateNewInstitution(e echo.Context) error
	UpdateInstitution(e echo.Context) error
	DeleteInstitution(e echo.Context) error
}

type InstitutionService struct {
	uc controller.InterfaceInstitutionController
}

func NewInstitutionService(uc controller.InterfaceInstitutionController) *InstitutionService {
	return &InstitutionService{uc: uc}
}

func (c *InstitutionService) GetAllInstitution(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllInstitution")
	defer span.Finish()

	res, err := c.uc.GetAllInstitution(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All Institution",
		Data:    res,
	})
}

func (c *InstitutionService) GetInstitutionByID(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetInstitutionByID")
	defer span.Finish()

	id := e.Param("id")
	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	res, err := c.uc.GetInstitutionByID(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Institution By ID",
		Data:    res,
	})
}

func (c *InstitutionService) CreateNewInstitution(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewInstitution")
	defer span.Finish()

	var institution *model.Institution
	if err := e.Bind(&institution); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	err := c.uc.InsertNewInstitution(ctx, institution)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", institution)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Institution",
		Data:    nil,
	})
}

func (c *InstitutionService) UpdateInstitution(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateInstitution")
	defer span.Finish()

	var institution *model.Institution
	if err := e.Bind(&institution); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	err := c.uc.UpdateInstitution(ctx, institution)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", institution)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Institution",
		Data:    nil,
	})
}

func (c *InstitutionService) DeleteInstitution(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteInstitution")
	defer span.Finish()

	id := e.Param("id")
	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	err := c.uc.DeleteInstitution(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Delete Institution")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete Institution",
		Data:    nil,
	})
}
