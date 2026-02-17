package domain

import "context"

type Presence struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	InstitutionID string `json:"institution_id"`
	Status        string `json:"status"`
	Note          string `json:"note"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type PresenceFilter struct {
	UserID        string
	InstitutionID string
	Status        string
}

type InterfacePresenceRepository interface {
	Create(ctx context.Context, req *Presence) (*Presence, error)
	Get(ctx context.Context, id string) (*Presence, error)
	List(ctx context.Context, filter PresenceFilter) ([]*Presence, error)
	Update(ctx context.Context, id string, req *Presence) (*Presence, error)
	Delete(ctx context.Context, id string) error
}
