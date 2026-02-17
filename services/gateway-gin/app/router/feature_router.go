package router

import "github.com/gin-gonic/gin"

func InitFeatureRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.feature

	route.GET("", RequirePermission("gateway.feature.read"), wrapHandler(service.GetAllFeatures))
	route.POST("", RequirePermission("gateway.feature.create"), wrapHandler(service.CreateFeature))
	route.PUT("", RequirePermission("gateway.feature.update"), wrapHandler(service.UpdateFeature))

	route.POST("/institution", RequirePermission("gateway.feature.set_institution"), wrapHandler(service.SetInstitutionFeature))
	route.GET("/institution/:id", RequirePermission("gateway.feature.read"), wrapHandler(service.GetInstitutionFeatures))
}
