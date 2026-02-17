package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"github.com/adrputra/face-recognition-svc/gateway/app/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type InterfaceDatasetClient = legacy.InterfaceDatasetClient
type DatasetClient = legacy.DatasetClient

func NewDatasetClient(db *gorm.DB, cfg *config.Config, mq *amqp.Channel) *DatasetClient {
	return legacy.NewDatasetClient(db, cfg, mq)
}
