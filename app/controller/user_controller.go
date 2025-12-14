package controller

import (
	"context"
	"encoding/json"
	"errors"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type InterfaceUserController interface {
	CreateNewUser(ctx context.Context, request *model.User) error
	GetUserDetail(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, request *model.User) error
	DeleteUser(ctx context.Context, username string) error
	Login(ctx context.Context, request *model.RequestLogin) (*model.ResponseLogin, error)
	GetAllUser(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error)
	GetInstitutionList(ctx context.Context) ([]string, error)

	UploadProfilePhoto(ctx context.Context, file *model.File) error
	UploadCoverPhoto(ctx context.Context, file *model.File) error
}

type UserController struct {
	userClient    client.InterfaceUserClient
	roleClient    client.InterfaceRoleClient
	paramClient   client.InterfaceParamClient
	storageClient client.InterfaceStorageClient
	config        *config.Config
	redis         *redis.Client
}

func NewUserController(userClient client.InterfaceUserClient, roleClient client.InterfaceRoleClient, paramClient client.InterfaceParamClient, storageClient client.InterfaceStorageClient, config *config.Config, redis *redis.Client) *UserController {
	return &UserController{
		userClient:    userClient,
		roleClient:    roleClient,
		paramClient:   paramClient,
		storageClient: storageClient,
		config:        config,
		redis:         redis,
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

	utils.LogEvent(span, "Username", username)

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Session", session)

	role, err := c.roleClient.GetRoleByID(ctx, session.RoleID)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}
	if role.Level == 3 {
		return nil, model.ThrowError(http.StatusUnauthorized, errors.New("you are not allowed to access this data (not authorized role)"))
	}

	user, err := c.userClient.GetUserDetail(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", user)

	if role.Level == 2 && user.InstitutionID != session.InstitutionID {
		return nil, model.ThrowError(http.StatusUnauthorized, errors.New("you are not allowed to access this data (different institution)"))
	}

	return user, nil
}

func (c *UserController) UpdateUser(ctx context.Context, request *model.User) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	err := c.userClient.UpdateUser(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}

func (c *UserController) DeleteUser(ctx context.Context, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	err := c.userClient.DeleteUser(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
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

	// Marshal menu mapping to JSON before storing in Redis
	menuMappingJSON, err := json.Marshal(menuMapping)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	ExpireCount, _ := strconv.Atoi(c.config.Auth.AccessExpiry)
	err = c.redis.Set(ctx, user.RoleID, menuMappingJSON, time.Duration(ExpireCount)*time.Hour).Err()
	if err != nil {
		utils.LogEventError(span, err)
		return nil, model.ThrowError(http.StatusInternalServerError, err)
	}

	response := &model.ResponseLogin{
		Username:        user.Username,
		Fullname:        user.Fullname,
		Shortname:       user.Shortname,
		Role:            user.RoleID,
		RoleName:        user.RoleName,
		Token:           accessToken,
		InstitutionID:   user.InstitutionID,
		InstitutionName: user.InstitutionName,
		MenuMapping:     role,
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (c *UserController) GetAllUser(ctx context.Context, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", pagination)
	utils.LogEvent(span, "Filter", filter)

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, nil, err
	}

	role, err := c.roleClient.GetRoleByID(ctx, session.RoleID)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, nil, err
	}

	users, pagination, err := c.userClient.GetAllUser(ctx, role.Level, session.InstitutionID, pagination, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, nil, err
	}

	utils.LogEvent(span, "Response", users)

	return users, pagination, nil
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

func (c *UserController) UploadProfilePhoto(ctx context.Context, file *model.File) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UploadProfilePhoto")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	res, err := c.storageClient.UploadFile(ctx, file, "bpkp", fmt.Sprintf("%s/%s", "profile-photo", session.Username))
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	err = c.userClient.UpdateProfilePhoto(ctx, res, session.Username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}

func (c *UserController) UploadCoverPhoto(ctx context.Context, file *model.File) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UploadCoverPhoto")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	res, err := c.storageClient.UploadFile(ctx, file, "bpkp", fmt.Sprintf("%s/%s", "cover-photo", session.Username))
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	err = c.userClient.UpdateCoverPhoto(ctx, res, session.Username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}
