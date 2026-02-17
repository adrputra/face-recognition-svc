package service

import (
	"errors"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/controller"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterfaceRoleService interface {
	CreateNewRoleMapping(e *gin.Context) error
	GetAllRoleMapping(e *gin.Context) error
	UpdateRoleMapping(e *gin.Context) error
	DeleteRoleMapping(e *gin.Context) error

	GetAllMenu(e *gin.Context) error
	CreateNewMenu(e *gin.Context) error
	UpdateMenu(e *gin.Context) error
	DeleteMenu(e *gin.Context) error

	GetAllRole(e *gin.Context) error
	CreateNewRole(e *gin.Context) error
}

type RoleService struct {
	uc controller.InterfaceRoleController
}

func NewRoleService(uc controller.InterfaceRoleController) InterfaceRoleService {
	return &RoleService{
		uc: uc,
	}
}

func (s *RoleService) CreateNewRoleMapping(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewRoleMapping")
	defer span.Finish()

	var request *model.MenuRoleMapping

	if err := e.ShouldBind(&request); err != nil {
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
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Role",
		Data:    nil,
	})
}

func (s *RoleService) GetAllRoleMapping(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllRoleMapping")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	response, pagination, err := s.uc.GetAllRoleMapping(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Role",
		Data:       response,
		Pagination: pagination,
	})
}

func (s *RoleService) UpdateRoleMapping(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateRoleMapping")
	defer span.Finish()

	var request *model.MenuRoleMapping

	if err := e.ShouldBind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.UpdateRoleMapping(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Update Role")
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Role",
		Data:    nil,
	})
}

func (s *RoleService) GetAllMenu(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllMenu")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	response, pagination, err := s.uc.GetAllMenu(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Menu",
		Data:       response,
		Pagination: pagination,
	})
}

func (s *RoleService) CreateNewMenu(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewMenu")
	defer span.Finish()

	var request *model.Menu

	if err := e.ShouldBind(&request); err != nil {
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
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Menu",
		Data:    nil,
	})
}

func (s *RoleService) GetAllRole(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllRole")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	response, pagination, err := s.uc.GetAllRole(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", response)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Role",
		Data:       response,
		Pagination: pagination,
	})
}

func (s *RoleService) CreateNewRole(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewRole")
	defer span.Finish()

	var request *model.Role

	if err := e.ShouldBind(&request); err != nil {
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
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New Role",
		Data:    nil,
	})
}

func (s *RoleService) UpdateMenu(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "UpdateRole")
	defer span.Finish()

	var request *model.Menu

	if err := e.ShouldBind(&request); err != nil {
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
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Menu",
		Data:    nil,
	})
}

func (s *RoleService) DeleteMenu(e *gin.Context) error {
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
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete Menu",
		Data:    nil,
	})
}

func (s *RoleService) DeleteRoleMapping(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteRoleMapping")
	defer span.Finish()

	id := e.Param("id")
	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	utils.LogEvent(span, "Request", id)

	err := s.uc.DeleteRoleMapping(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Delete Role Mapping")
	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Delete Role Mapping",
		Data:    nil,
	})
}
