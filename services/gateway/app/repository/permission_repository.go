package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"gorm.io/gorm"
)

type InterfacePermissionClient = legacy.InterfacePermissionClient
type PermissionClient = legacy.PermissionClient

func NewPermissionClient(db *gorm.DB) *PermissionClient {
	return legacy.NewPermissionClient(db)
}
