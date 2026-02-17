package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"gorm.io/gorm"
)

type InterfaceInstitutionClient = legacy.InterfaceInstitutionClient
type InstitutionClient = legacy.InstitutionClient

func NewInstitutionClient(db *gorm.DB) InterfaceInstitutionClient {
	return legacy.NewInstitutionClient(db)
}
