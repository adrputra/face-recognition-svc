package controller

import (
	"context"
	"errors"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type InterfaceUserController interface {
	CreateNewUser(ctx context.Context, request *model.User) error
	GetUserDetail(ctx context.Context, username string) (*model.User, error)
	Login(ctx context.Context, request *model.RequestLogin) (*model.ResponseLogin, error)
	GetAllUser(ctx context.Context) ([]*model.User, error)
	GetInstitutionList(ctx context.Context) ([]string, error)
}

type UserController struct {
	userClient client.InterfaceUserClient
	roleClient client.InterfaceRoleClient
}

func NewUserController(userClient client.InterfaceUserClient, roleClient client.InterfaceRoleClient) *UserController {
	return &UserController{
		userClient: userClient,
		roleClient: roleClient,
	}
}

func (c *UserController) CreateNewUser(ctx context.Context, request *model.User) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreateNewUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.Password = string(hashPassword)

	err = c.userClient.CreateNewUser(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}
	return nil
}

func (c *UserController) GetUserDetail(ctx context.Context, username string) (*model.User, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetUserDetail")
	defer span.Finish()

	span.LogKV("Request", username)

	user, err := c.userClient.GetUserDetail(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return user, nil
}

func (c *UserController) Login(ctx context.Context, request *model.RequestLogin) (*model.ResponseLogin, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: Login")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	user, err := c.userClient.GetUserDetail(ctx, request.Username)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		utils.LogEventError(span, errors.New("invalid username or password "))
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("invalid username or password "))
	}

	role, err := c.roleClient.GetMenuRoleMapping(ctx, user.RoleID)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	if len(role) < 1 {
		utils.LogEventError(span, errors.New("menu role mapping not found"))
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("menu role mapping not found"))
	}

	menuMapping := make(map[string]string)
	for _, v := range role {
		menuMapping[v.MenuID] = v.AccessMethod
	}

	accessToken, _, err := c.userClient.CreateAccessToken(ctx, user, false, menuMapping)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	response := &model.ResponseLogin{
		Username:      user.Username,
		Fullname:      user.Fullname,
		Shortname:     user.Shortname,
		Role:          user.RoleID,
		Token:         accessToken,
		InstitutionID: user.InstitutionID,
		MenuMapping:   role,
	}

	return response, nil
}

func (c *UserController) GetAllUser(ctx context.Context) ([]*model.User, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllUser")
	defer span.Finish()

	users, err := c.userClient.GetAllUser(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", users)

	return users, nil
}

func (c *UserController) GetInstitutionList(ctx context.Context) ([]string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetInstitutionList")
	defer span.Finish()

	institutionList, err := c.userClient.GetInstitutionList(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", institutionList)

	return institutionList, nil
}
