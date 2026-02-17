package router

import (
	"github.com/adrputra/face-recognition-svc/gateway/app/config"
	"github.com/adrputra/face-recognition-svc/gateway/app/controller"
	"github.com/adrputra/face-recognition-svc/gateway/app/repository"
	"github.com/adrputra/face-recognition-svc/gateway/app/service"
	"github.com/adrputra/face-recognition-svc/gateway/app/utils"

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
	permission  service.InterfacePermissionService
	feature     service.InterfaceFeatureService
	presence    service.InterfacePresenceService
}

type ControllerFactory struct {
	user        controller.InterfaceUserController
	dataset     controller.InterfaceDatasetController
	role        controller.InterfaceRoleController
	param       controller.InterfaceParamController
	institution controller.InterfaceInstitutionController
	permission  controller.InterfacePermissionController
	feature     controller.InterfaceFeatureController
	presence    controller.InterfacePresenceController
}

type RepositoryFactory struct {
	user        repository.InterfaceUserClient
	storage     repository.InterfaceStorageClient
	role        repository.InterfaceRoleClient
	permission  repository.InterfacePermissionClient
	feature     repository.InterfaceFeatureClient
	dataset     repository.InterfaceDatasetClient
	param       repository.InterfaceParamClient
	institution repository.InterfaceInstitutionClient
}

type MiddlewareFactory struct {
	Auth utils.InterfaceAuthMiddleware
}

type Factory struct {
	Service    ServiceFactory
	Controller ControllerFactory
	Repository RepositoryFactory
	Middleware MiddlewareFactory
}

var factory *Factory

func InitFactory(cfg *config.Config, db *gorm.DB, s3 *s3.S3, redis *redis.Client, mq *amqp.Channel) {
	repositoryFactory := RepositoryFactory{
		user:        repository.NewUserClient(db, cfg),
		storage:     repository.NewStorageClient(s3, db),
		role:        repository.NewRoleClient(db),
		permission:  repository.NewPermissionClient(db),
		feature:     repository.NewFeatureClient(db),
		dataset:     repository.NewDatasetClient(db, cfg, mq),
		param:       repository.NewParamClient(db),
		institution: repository.NewInstitutionClient(db),
	}
	controller := ControllerFactory{
		user:        controller.NewUserController(repositoryFactory.user, repositoryFactory.role, repositoryFactory.param, repositoryFactory.storage, cfg, redis),
		dataset:     controller.NewDatasetController(repositoryFactory.storage, db, repositoryFactory.user, cfg, repositoryFactory.dataset),
		role:        controller.NewRoleController(repositoryFactory.role),
		permission:  controller.NewPermissionController(repositoryFactory.permission),
		feature:     controller.NewFeatureController(repositoryFactory.feature),
		param:       controller.NewParamController(redis, repositoryFactory.param),
		institution: controller.NewInstitutionController(repositoryFactory.institution),
		presence:    controller.NewPresenceController(repository.NewPresenceRepository(cfg)),
	}
	service := ServiceFactory{
		user:        service.NewUserService(controller.user),
		dataset:     service.NewDatasetService(controller.dataset),
		role:        service.NewRoleService(controller.role),
		permission:  service.NewPermissionService(controller.permission),
		feature:     service.NewFeatureService(controller.feature),
		param:       service.NewParamService(controller.param),
		institution: service.NewInstitutionService(controller.institution),
		presence:    service.NewPresenceService(controller.presence),
	}

	middleware := MiddlewareFactory{
		Auth: utils.NewAuthMiddleware(db, redis),
	}
	factory = &Factory{
		Service:    service,
		Controller: controller,
		Repository: repositoryFactory,
		Middleware: middleware,
	}
}

func GetFactory() *Factory {
	return factory
}
