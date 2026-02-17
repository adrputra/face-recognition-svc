package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"gorm.io/gorm"
)

type InterfaceRoleClient = legacy.InterfaceRoleClient
type RoleClient = legacy.RoleClient

func NewRoleClient(db *gorm.DB) *RoleClient {
	return legacy.NewRoleClient(db)
}
