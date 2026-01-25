package service

import (
	"errors"
	"face-recognition-svc/gateway/app/controller"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfacePermissionService interface {
	GetAllPermissions(e echo.Context) error
	CreatePermission(e echo.Context) error
	UpdatePermission(e echo.Context) error
	AssignRolePermissions(e echo.Context) error
	GetRolePermissions(e echo.Context) error
}

type PermissionService struct {
	pc controller.InterfacePermissionController
}

func NewPermissionService(pc controller.InterfacePermissionController) InterfacePermissionService {
	return &PermissionService{pc: pc}
}

func (s *PermissionService) GetAllPermissions(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAllPermissions")
	defer span.Finish()

	pagination := utils.ParsePaginationFromQuery(e)
	filter := utils.ParseFilterFromQuery(e)

	res, pagination, err := s.pc.GetAllPermissions(ctx, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:       200,
		Message:    "Success Get All Permission",
		Data:       res,
		Pagination: pagination,
	})
}

func (s *PermissionService) CreatePermission(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreatePermission")
	defer span.Finish()

	var request *model.Permission
	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.pc.CreatePermission(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create Permission",
		Data:    nil,
	})
}

func (s *PermissionService) UpdatePermission(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdatePermission")
	defer span.Finish()

	var request *model.PermissionUpdateRequest
	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.pc.UpdatePermission(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Update Permission",
		Data:    nil,
	})
}

func (s *PermissionService) AssignRolePermissions(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "AssignRolePermissions")
	defer span.Finish()

	var request *model.RolePermissionAssignment
	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	if request == nil {
		return utils.LogError(e, errors.New("invalid request"), nil)
	}

	if err := s.pc.AssignRolePermissions(ctx, request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Assign Role Permissions",
		Data:    nil,
	})
}

func (s *PermissionService) GetRolePermissions(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetRolePermissions")
	defer span.Finish()

	roleID := e.Param("id")
	if roleID == "" {
		return utils.LogError(e, errors.New("role_id shouldn't be empty"), nil)
	}

	res, err := s.pc.GetRolePermissions(ctx, roleID)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Role Permissions",
		Data:    res,
	})
}
