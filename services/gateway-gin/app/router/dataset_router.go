package router

import "github.com/gin-gonic/gin"

func InitDatasetRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.dataset

	route.GET("", RequirePermission("gateway.dataset.read"), wrapHandler(service.GetDatasetList))
	route.POST("", RequirePermission("gateway.dataset.create"), wrapHandler(service.UploadUserDataset))
	route.DELETE("/:id", RequirePermission("gateway.dataset.delete"), wrapHandler(service.DeleteDataset))

	route.POST("/train-model/:id", RequirePermission("gateway.dataset.train_model"), wrapHandler(service.TrainModel))
	route.GET("/last-train-model/:id", RequirePermission("gateway.dataset.read"), wrapHandler(service.GetLastTrainModel))

	// TODO: POST is used for a read-like query; confirm action naming.
	route.POST("/model-training-history", RequirePermission("gateway.dataset.model_training_history"), wrapHandler(service.GetModelTrainingHistory))

	route.GET("/:institution-id/:id", RequirePermission("gateway.dataset.read"), wrapHandler(service.GetDatasetsByUsername))
}
