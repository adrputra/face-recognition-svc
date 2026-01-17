package router

import "github.com/labstack/echo/v4"

func InitDatasetRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.dataset

	route.GET("", service.GetDatasetList)
	route.POST("", service.UploadUserDataset)
	route.DELETE("/:id", service.DeleteDataset)

	route.POST("/train-model/:id", service.TrainModel)
	route.GET("/last-train-model/:id", service.GetLastTrainModel)

	route.POST("/model-training-history", service.GetModelTrainingHistory)

	route.GET("/:institution-id/:id", service.GetDatasetsByUsername)
}
