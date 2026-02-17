package domain

import "time"

type Entity struct {
	ID            string
	UserID        string
	InstitutionID string
	Status        string
	Note          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Filter struct {
	UserID        string
	InstitutionID string
	Status        string
}

type Patch struct {
	UserID        *string
	InstitutionID *string
	Status        *string
	Note          *string
}
