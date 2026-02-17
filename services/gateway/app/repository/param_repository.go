package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"gorm.io/gorm"
)

type InterfaceParamClient = legacy.InterfaceParamClient
type ParamClient = legacy.ParamClient

func NewParamClient(db *gorm.DB) *ParamClient {
	return legacy.NewParamClient(db)
}
