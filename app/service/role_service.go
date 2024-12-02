package service

import (
	"errors"
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceRoleService interface {
	CreateNewRoleMapping(e echo.Context) error
	GetAllRoleMapping(e echo.Context) error

	GetAllMenu(e echo.Context) error
	CreateNewMenu(e echo.Context) error
	UpdateMenu(e echo.Context) error
	DeleteMenu(e echo.Context) error

	GetAllRole(e echo.Context) error
	CreateNewRole(e echo.Context) error
}

type RoleService struct {
	uc controller.InterfaceRoleController
}

func NewRoleService(uc controller.InterfaceRoleController) InterfaceRoleService {
	return &RoleService{
		uc: uc,
	}
}

func (s *RoleService) CreateNewRoleMapping(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewRoleMapping")
	defer span.Finish()

	var request *model.MenuRoleMapping

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CreateNewRoleMapping(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Create New Role")
	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Role",
		Data:    nil,
	})
}

func (s *RoleService) GetAllRoleMapping(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllRoleMapping")
	defer span.Finish()

	response, err := s.uc.GetAllRoleMapping(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All Role",
		Data:    response,
	})
}

func (s *RoleService) GetAllMenu(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllMenu")
	defer span.Finish()

	response, err := s.uc.GetAllMenu(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All Menu",
		Data:    response,
	})
}

func (s *RoleService) CreateNewMenu(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewMenu")
	defer span.Finish()

	var request *model.Menu

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CreateNewMenu(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Create New Menu")
	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Menu",
		Data:    nil,
	})
}

func (s *RoleService) GetAllRole(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllRole")
	defer span.Finish()

	response, err := s.uc.GetAllRole(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All Role",
		Data:    response,
	})
}

func (s *RoleService) CreateNewRole(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewRole")
	defer span.Finish()

	var request *model.Role

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CreateNewRole(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Create New Role")
	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Role",
		Data:    nil,
	})
}

func (s *RoleService) UpdateMenu(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateRole")
	defer span.Finish()

	var request *model.Menu

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.UpdateMenu(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Update Menu")
	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Menu",
		Data:    nil,
	})
}

func (s *RoleService) DeleteMenu(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteMenu")
	defer span.Finish()

	id := e.Param("id")
	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	utils.LogEvent(span, "Request", id)

	err := s.uc.DeleteMenu(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Delete Menu")
	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete Menu",
		Data:    nil,
	})
}
