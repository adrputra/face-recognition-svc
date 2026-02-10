package router

import "github.com/labstack/echo/v4"

func InitFeatureRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.feature

	route.GET("", service.GetAllFeatures, RequirePermission("gateway.feature.read"))
	route.POST("", service.CreateFeature, RequirePermission("gateway.feature.create"))
	route.PUT("", service.UpdateFeature, RequirePermission("gateway.feature.update"))

	route.POST("/institution", service.SetInstitutionFeature, RequirePermission("gateway.feature.set_institution"))
	route.GET("/institution/:id", service.GetInstitutionFeatures, RequirePermission("gateway.feature.read"))
}
