package controller

import (
	"context"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"time"

	"github.com/google/uuid"
)

type InterfaceRoleController interface {
	CreateNewRoleMapping(ctx context.Context, request *model.MenuRoleMapping) error
	GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error)

	GetAllMenu(ctx context.Context) ([]*model.Menu, error)
	CreateNewMenu(ctx context.Context, request *model.Menu) error
	UpdateMenu(ctx context.Context, request *model.Menu) error
	DeleteMenu(ctx context.Context, id string) error

	GetAllRole(ctx context.Context) ([]*model.Role, error)
	CreateNewRole(ctx context.Context, request *model.Role) error
}

type RoleController struct {
	roleClient client.InterfaceRoleClient
}

func NewRoleController(roleClient client.InterfaceRoleClient) *RoleController {
	return &RoleController{
		roleClient: roleClient,
	}
}

func (c *RoleController) CreateNewRoleMapping(ctx context.Context, request *model.MenuRoleMapping) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreateNewRoleMapping")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.CreatedBy = session.Username
	request.UpdatedBy = session.Username

	err = c.roleClient.CreateNewRoleMapping(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}
	return nil
}

func (c *RoleController) GetAllRoleMapping(ctx context.Context) ([]*model.MenuRoleMapping, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllRoleMapping")
	defer span.Finish()

	response, err := c.roleClient.GetAllRoleMapping(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (c *RoleController) GetAllMenu(ctx context.Context) ([]*model.Menu, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllMenu")
	defer span.Finish()

	response, err := c.roleClient.GetAllMenu(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (c *RoleController) CreateNewMenu(ctx context.Context, request *model.Menu) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreateNewMenu")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.CreatedBy = session.Username
	request.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", request)

	err = c.roleClient.CreateNewMenu(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}
	return nil
}

func (c *RoleController) GetAllRole(ctx context.Context) ([]*model.Role, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllRole")
	defer span.Finish()

	response, err := c.roleClient.GetAllRole(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (c *RoleController) CreateNewRole(ctx context.Context, request *model.Role) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CreateNewRole")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.Id = uuid.New().String()
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.CreatedBy = session.Username
	request.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", request)

	err = c.roleClient.CreateNewRole(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}
	return nil
}

func (c *RoleController) UpdateMenu(ctx context.Context, request *model.Menu) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateMenu")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.UpdatedAt = time.Now()
	request.UpdatedBy = session.Username

	utils.LogEvent(span, "Request", request)

	err = c.roleClient.UpdateMenu(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Menu")

	return nil
}

func (c *RoleController) DeleteMenu(ctx context.Context, id string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteMenu")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)

	err := c.roleClient.DeleteMenu(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Menu")

	return nil
}
