package router

import "github.com/gin-gonic/gin"

func InitParamRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.param

	route.GET("/:id", RequirePermission("gateway.param.read"), wrapHandler(service.GetParameterByKey))
	route.GET("", RequirePermission("gateway.param.read"), wrapHandler(service.GetAllParam))
	route.POST("", RequirePermission("gateway.param.create"), wrapHandler(service.InsertNewParam))
	route.PUT("", RequirePermission("gateway.param.update"), wrapHandler(service.UpdateParam))
	route.DELETE("/:id", RequirePermission("gateway.param.delete"), wrapHandler(service.DeleteParam))
}
