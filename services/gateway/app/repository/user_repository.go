package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"github.com/adrputra/face-recognition-svc/gateway/app/config"
	"gorm.io/gorm"
)

type InterfaceUserClient = legacy.InterfaceUserClient
type UserClient = legacy.UserClient

func NewUserClient(db *gorm.DB, cfg *config.Config) *UserClient {
	return legacy.NewUserClient(db, cfg)
}
