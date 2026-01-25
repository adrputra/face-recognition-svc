package router

import "github.com/labstack/echo/v4"

func InitFeatureRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.feature

	route.GET("", service.GetAllFeatures)
	route.POST("", service.CreateFeature)
	route.PUT("", service.UpdateFeature)

	route.POST("/institution", service.SetInstitutionFeature)
	route.GET("/institution/:id", service.GetInstitutionFeatures)
}
