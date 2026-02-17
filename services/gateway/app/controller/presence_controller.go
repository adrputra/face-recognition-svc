package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/adrputra/face-recognition-svc/gateway/app/domain"
)

var (
	ErrPresenceRequiredID          = errors.New("id shouldn't be empty")
	ErrPresenceRequiredUserID      = errors.New("user_id is required")
	ErrPresenceRequiredUpdateField = errors.New("at least one field must be provided")
)

type InterfacePresenceController interface {
	CreatePresence(ctx context.Context, req *domain.Presence) (*domain.Presence, error)
	GetPresence(ctx context.Context, id string) (*domain.Presence, error)
	ListPresence(ctx context.Context, filter domain.PresenceFilter) ([]*domain.Presence, error)
	UpdatePresence(ctx context.Context, id string, req *domain.Presence) (*domain.Presence, error)
	DeletePresence(ctx context.Context, id string) error
}

type PresenceController struct {
	repository domain.InterfacePresenceRepository
}

func NewPresenceController(repository domain.InterfacePresenceRepository) InterfacePresenceController {
	return &PresenceController{
		repository: repository,
	}
}

func (c *PresenceController) CreatePresence(ctx context.Context, req *domain.Presence) (*domain.Presence, error) {
	if req == nil || strings.TrimSpace(req.UserID) == "" {
		return nil, ErrPresenceRequiredUserID
	}

	req.UserID = strings.TrimSpace(req.UserID)
	req.InstitutionID = strings.TrimSpace(req.InstitutionID)
	req.Status = strings.TrimSpace(req.Status)
	return c.repository.Create(ctx, req)
}

func (c *PresenceController) GetPresence(ctx context.Context, id string) (*domain.Presence, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrPresenceRequiredID
	}
	return c.repository.Get(ctx, id)
}

func (c *PresenceController) ListPresence(ctx context.Context, filter domain.PresenceFilter) ([]*domain.Presence, error) {
	filter.UserID = strings.TrimSpace(filter.UserID)
	filter.InstitutionID = strings.TrimSpace(filter.InstitutionID)
	filter.Status = strings.TrimSpace(filter.Status)
	return c.repository.List(ctx, filter)
}

func (c *PresenceController) UpdatePresence(ctx context.Context, id string, req *domain.Presence) (*domain.Presence, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrPresenceRequiredID
	}
	if req == nil {
		return nil, ErrPresenceRequiredUpdateField
	}

	req.UserID = strings.TrimSpace(req.UserID)
	req.InstitutionID = strings.TrimSpace(req.InstitutionID)
	req.Status = strings.TrimSpace(req.Status)
	req.Note = strings.TrimSpace(req.Note)
	if req.UserID == "" && req.InstitutionID == "" && req.Status == "" && req.Note == "" {
		return nil, ErrPresenceRequiredUpdateField
	}

	return c.repository.Update(ctx, id, req)
}

func (c *PresenceController) DeletePresence(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrPresenceRequiredID
	}
	return c.repository.Delete(ctx, id)
}
