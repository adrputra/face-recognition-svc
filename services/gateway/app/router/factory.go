package router

import (
	"face-recognition-svc/gateway/app/client"
	"face-recognition-svc/gateway/app/config"
	"face-recognition-svc/gateway/app/controller"
	"face-recognition-svc/gateway/app/service"
	"face-recognition-svc/gateway/app/utils"

	"github.com/aws/aws-sdk-go/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	user        service.InterfaceUserService
	dataset     service.InterfaceDatasetService
	role        service.InterfaceRoleService
	param       service.InterfaceParamService
	institution service.InterfaceInstitutionService
}

type ControllerFactory struct {
	user        controller.InterfaceUserController
	dataset     controller.InterfaceDatasetController
	role        controller.InterfaceRoleController
	param       controller.InterfaceParamController
	institution controller.InterfaceInstitutionController
}

type ClientFactory struct {
	user        client.InterfaceUserClient
	storage     client.InterfaceStorageClient
	role        client.InterfaceRoleClient
	dataset     client.InterfaceDatasetClient
	param       client.InterfaceParamClient
	institution client.InterfaceInstitutionClient
}

type MiddlewareFactory struct {
	Auth utils.InterfaceAuthMiddleware
}

type Factory struct {
	Service    ServiceFactory
	Controller ControllerFactory
	Client     ClientFactory
	Middleware MiddlewareFactory
}

var factory *Factory

func InitFactory(cfg *config.Config, db *gorm.DB, s3 *s3.S3, redis *redis.Client, mq *amqp.Channel) {
	client := ClientFactory{
		user:        client.NewUserClient(db, cfg),
		storage:     client.NewStorageClient(s3, db),
		role:        client.NewRoleClient(db),
		dataset:     client.NewDatasetClient(db, cfg, mq),
		param:       client.NewParamClient(db),
		institution: client.NewInstitutionClient(db),
	}
	controller := ControllerFactory{
		user:        controller.NewUserController(client.user, client.role, client.param, client.storage, cfg, redis),
		dataset:     controller.NewDatasetController(client.storage, db, client.user, cfg, client.dataset),
		role:        controller.NewRoleController(client.role),
		param:       controller.NewParamController(redis, client.param),
		institution: controller.NewInstitutionController(client.institution),
	}
	service := ServiceFactory{
		user:        service.NewUserService(controller.user),
		dataset:     service.NewDatasetService(controller.dataset),
		role:        service.NewRoleService(controller.role),
		param:       service.NewParamService(controller.param),
		institution: service.NewInstitutionService(controller.institution),
	}
	middleware := MiddlewareFactory{
		Auth: utils.NewAuthMiddleware(db, redis),
	}
	factory = &Factory{
		Service:    service,
		Controller: controller,
		Client:     client,
		Middleware: middleware,
	}
}

func GetFactory() *Factory {
	return factory
}
