package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"gorm.io/gorm"
)

type InterfaceFeatureClient = legacy.InterfaceFeatureClient
type FeatureClient = legacy.FeatureClient

func NewFeatureClient(db *gorm.DB) *FeatureClient {
	return legacy.NewFeatureClient(db)
}
