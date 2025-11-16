package client

import (
	"context"
	"errors"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"gorm.io/gorm"
)

type InterfaceInstitutionClient interface {
	GetAllInstitutions(ctx context.Context) ([]*model.Institution, error)
	GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error)
	CreateNewInstitution(ctx context.Context, institution *model.Institution) error
	UpdateInstitution(ctx context.Context, institution *model.Institution) error
	DeleteInstitution(ctx context.Context, id string) error
}

type InstitutionClient struct {
	db *gorm.DB
}

func NewInstitutionClient(db *gorm.DB) InterfaceInstitutionClient {
	return &InstitutionClient{db: db}
}

func (c *InstitutionClient) GetAllInstitutions(ctx context.Context) ([]*model.Institution, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetAllInstitutions")
	defer span.Finish()

	utils.LogEvent(span, "Request", "All")
	var response []*model.Institution

	query := "SELECT * FROM institution"
	utils.LogEvent(span, "Query", query)

	err := c.db.Debug().Raw(query).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	if response == nil {
		return nil, model.ThrowError(http.StatusInternalServerError, errors.New("Data Not Found"))
	}

	utils.LogEvent(span, "Response", response)
	return response, nil
}

func (c *InstitutionClient) GetInstitutionByID(ctx context.Context, id string) (*model.Institution, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetInstitutionByID")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)
	var response *model.Institution

	query := "SELECT * FROM institution WHERE id = ?"
	utils.LogEvent(span, "Query", query)

	err := c.db.Debug().Raw(query, id).Scan(&response).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", response)
	return response, nil
}

func (c *InstitutionClient) CreateNewInstitution(ctx context.Context, institution *model.Institution) error {
	span, _ := utils.SpanFromContext(ctx, "Client: CreateNewInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", institution)

	var args []interface{}

	args = append(args, institution.ID, institution.Name, institution.Address, institution.PhoneNumber, institution.Email, institution.CreatedAt, institution.CreatedBy)
	err := c.db.Debug().Exec("INSERT INTO institution (id, name, address, phone_number, email, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?)", args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Create New Institution")
	return nil
}

func (c *InstitutionClient) UpdateInstitution(ctx context.Context, institution *model.Institution) error {
	span, _ := utils.SpanFromContext(ctx, "Client: UpdateInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", institution)

	var args []interface{}
	args = append(args, institution.Name, institution.Address, institution.PhoneNumber, institution.UpdatedAt, institution.UpdatedBy, institution.ID)
	err := c.db.Debug().Exec("UPDATE institution SET name = ?, address = ?, phone_number = ?, updated_at = ?, updated_by = ? WHERE id = ?", args...).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Update Institution")
	return nil
}

func (c *InstitutionClient) DeleteInstitution(ctx context.Context, id string) error {
	span, _ := utils.SpanFromContext(ctx, "Client: DeleteInstitution")
	defer span.Finish()

	utils.LogEvent(span, "Request", id)

	err := c.db.Debug().Exec("DELETE FROM institution WHERE id = ?", id).Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Institution")
	return nil
}
