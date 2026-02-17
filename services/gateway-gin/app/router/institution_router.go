package router

import "github.com/gin-gonic/gin"

func InitInstitutionRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.institution

	route.GET("", RequirePermission("gateway.institution.read"), wrapHandler(service.GetAllInstitution))
	route.GET("/:id", RequirePermission("gateway.institution.read"), wrapHandler(service.GetInstitutionByID))
	route.POST("", RequirePermission("gateway.institution.create"), wrapHandler(service.CreateNewInstitution))
	route.PUT("", RequirePermission("gateway.institution.update"), wrapHandler(service.UpdateInstitution))
	route.DELETE("/:id", RequirePermission("gateway.institution.delete"), wrapHandler(service.DeleteInstitution))
}
