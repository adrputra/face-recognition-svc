package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	domain "github.com/adrputra/face-recognition-svc/presence-svc/internal/domain"
)

var (
	ErrRequiredUserID      = errors.New("user_id is required")
	ErrRequiredID          = errors.New("id is required")
	ErrRequiredUpdateField = errors.New("at least one field must be provided")
	ErrEmptyUserID         = errors.New("user_id cannot be empty")
	ErrEmptyStatus         = errors.New("status cannot be empty")
)

type InterfacePresenceController interface {
	CreatePresence(ctx context.Context, req *domain.Entity) (*domain.Entity, error)
	GetPresence(ctx context.Context, id string) (*domain.Entity, error)
	ListPresence(ctx context.Context, filter domain.Filter) ([]*domain.Entity, error)
	UpdatePresence(ctx context.Context, id string, patch domain.Patch) (*domain.Entity, error)
	DeletePresence(ctx context.Context, id string) error
}

type PresenceController struct {
	repository domain.Repository
}

func NewPresenceController(repository domain.Repository) InterfacePresenceController {
	return &PresenceController{
		repository: repository,
	}
}

func (c *PresenceController) CreatePresence(ctx context.Context, req *domain.Entity) (*domain.Entity, error) {
	if req == nil {
		return nil, ErrRequiredUserID
	}

	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		return nil, ErrRequiredUserID
	}

	statusValue := strings.TrimSpace(req.Status)
	if statusValue == "" {
		statusValue = "present"
	}

	now := time.Now().UTC()
	item := &domain.Entity{
		ID:            newID(),
		UserID:        userID,
		InstitutionID: strings.TrimSpace(req.InstitutionID),
		Status:        statusValue,
		Note:          req.Note,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return c.repository.Create(ctx, item)
}

func (c *PresenceController) GetPresence(ctx context.Context, id string) (*domain.Entity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrRequiredID
	}
	return c.repository.Get(ctx, id)
}

func (c *PresenceController) ListPresence(ctx context.Context, filter domain.Filter) ([]*domain.Entity, error) {
	filter.UserID = strings.TrimSpace(filter.UserID)
	filter.InstitutionID = strings.TrimSpace(filter.InstitutionID)
	filter.Status = strings.TrimSpace(filter.Status)
	return c.repository.List(ctx, filter)
}

func (c *PresenceController) UpdatePresence(ctx context.Context, id string, patch domain.Patch) (*domain.Entity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrRequiredID
	}

	hasField := false
	if patch.UserID != nil {
		userID := strings.TrimSpace(*patch.UserID)
		if userID == "" {
			return nil, ErrEmptyUserID
		}
		patch.UserID = &userID
		hasField = true
	}
	if patch.InstitutionID != nil {
		institutionID := strings.TrimSpace(*patch.InstitutionID)
		patch.InstitutionID = &institutionID
		hasField = true
	}
	if patch.Status != nil {
		statusValue := strings.TrimSpace(*patch.Status)
		if statusValue == "" {
			return nil, ErrEmptyStatus
		}
		patch.Status = &statusValue
		hasField = true
	}
	if patch.Note != nil {
		hasField = true
	}

	if !hasField {
		return nil, ErrRequiredUpdateField
	}

	return c.repository.Update(ctx, id, patch)
}

func (c *PresenceController) DeletePresence(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrRequiredID
	}
	return c.repository.Delete(ctx, id)
}

func newID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return strconvFallbackID()
	}
	return hex.EncodeToString(raw)
}

func strconvFallbackID() string {
	return strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339Nano), ":", "")
}
