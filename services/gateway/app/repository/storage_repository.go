package repository

import (
	legacy "github.com/adrputra/face-recognition-svc/gateway/app/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
)

type InterfaceStorageClient = legacy.InterfaceStorageClient
type StorageClient = legacy.StorageClient

func NewStorageClient(s3Client *s3.S3, db *gorm.DB) *StorageClient {
	return legacy.NewStorageClient(s3Client, db)
}
