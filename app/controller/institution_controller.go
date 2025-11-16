package controller

import (
	"context"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"

	"github.com/google/uuid"
)

type InterfaceInstitutionController interface {
	GetAllInstitution(ctx context.Context) ([]*model.Institution, error)
	GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error)
	InsertNewInstitution(ctx context.Context, institution *model.Institution) error
	UpdateInstitution(ctx context.Context, institution *model.Institution) error
	DeleteInstitution(ctx context.Context, id string) error
}

type InstitutionController struct {
	institutionClient client.InterfaceInstitutionClient
}

func NewInstitutionController(institutionClient client.InterfaceInstitutionClient) *InstitutionController {
	return &InstitutionController{institutionClient: institutionClient}
}

func (c *InstitutionController) GetAllInstitution(ctx context.Context) ([]*model.Institution, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetAllInstitution")
	defer span.Finish()

	res, err := c.institutionClient.GetAllInstitutions(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return res, nil
}

func (c *InstitutionController) GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetInstitutionByID")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)
	res, err := c.institutionClient.GetInstitutionByID(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)
	return res, nil
}

func (c *InstitutionController) InsertNewInstitution(ctx context.Context, institution *model.Institution) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: InsertNewInstitution")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}
	institution.ID = uuid.New().String()
	institution.CreatedAt = utils.LocalTime().Format("2006-01-02 15:04:05")
	institution.CreatedBy = session.Username
	utils.LogEvent(span, "Request", institution)

	err = c.institutionClient.CreateNewInstitution(ctx, institution)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Insert New Institution")
	return nil
}

func (c *InstitutionController) UpdateInstitution(ctx context.Context, institution *model.Institution) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UpdateInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", institution)
	err := c.institutionClient.UpdateInstitution(ctx, institution)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Institution")
	return nil
}

func (c *InstitutionController) DeleteInstitution(ctx context.Context, id string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)
	err := c.institutionClient.DeleteInstitution(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Institution")
	return nil
}
