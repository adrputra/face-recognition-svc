package service

import (
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceUserService interface {
	CreateNewUser(e echo.Context) error
	GetUserDetail(e echo.Context) error
	Login(e echo.Context) error
	GetAllUser(e echo.Context) error
	GetInstitutionList(e echo.Context) error
}

type UserService struct {
	uc controller.InterfaceUserController
}

func NewUserService(uc controller.InterfaceUserController) InterfaceUserService {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) CreateNewUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreateNewUser")
	defer span.Finish()

	var request *model.User

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CreateNewUser(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Success Create New User")

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Create New User",
		Data:    nil,
	})
}

func (s *UserService) GetUserDetail(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetUserDetail")
	defer span.Finish()

	username := e.Param("id")

	utils.LogEvent(span, "Request", username)

	user, err := s.uc.GetUserDetail(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get User Detail",
		Data:    user,
	})
}

func (s *UserService) Login(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "Login")
	defer span.Finish()

	var request *model.RequestLogin

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	response, err := s.uc.Login(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Login",
		Data:    response,
	})
}

func (s *UserService) GetAllUser(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetAlluser")
	defer span.Finish()

	users, err := s.uc.GetAllUser(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", users)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get All User",
		Data:    users,
	})
}

func (s *UserService) GetInstitutionList(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetInstitutionList")
	defer span.Finish()

	institutionList, err := s.uc.GetInstitutionList(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", institutionList)

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Institution List",
		Data:    institutionList,
	})
}
