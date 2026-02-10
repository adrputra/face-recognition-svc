package router

import "github.com/labstack/echo/v4"

func InitDatasetRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.dataset

	route.GET("", service.GetDatasetList, RequirePermission("gateway.dataset.read"))
	route.POST("", service.UploadUserDataset, RequirePermission("gateway.dataset.create"))
	route.DELETE("/:id", service.DeleteDataset, RequirePermission("gateway.dataset.delete"))

	route.POST("/train-model/:id", service.TrainModel, RequirePermission("gateway.dataset.train_model"))
	route.GET("/last-train-model/:id", service.GetLastTrainModel, RequirePermission("gateway.dataset.read"))

	// TODO: POST is used for a read-like query; confirm action naming.
	route.POST("/model-training-history", service.GetModelTrainingHistory, RequirePermission("gateway.dataset.model_training_history"))

	route.GET("/:institution-id/:id", service.GetDatasetsByUsername, RequirePermission("gateway.dataset.read"))
}
