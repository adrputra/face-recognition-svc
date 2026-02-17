package router

import "github.com/gin-gonic/gin"

func InitRoleRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.role

	route.GET("", RequirePermission("gateway.role.read"), wrapHandler(service.GetAllRole))

	route.GET("/mapping", RequirePermission("gateway.role_mapping.read"), wrapHandler(service.GetAllRoleMapping))
	route.POST("/create", RequirePermission("gateway.role.create"), wrapHandler(service.CreateNewRole))
	route.POST("/mapping/create", RequirePermission("gateway.role_mapping.create"), wrapHandler(service.CreateNewRoleMapping))
	route.PUT("/mapping", RequirePermission("gateway.role_mapping.update"), wrapHandler(service.UpdateRoleMapping))
	route.DELETE("/mapping/:id", RequirePermission("gateway.role_mapping.delete"), wrapHandler(service.DeleteRoleMapping))

	route.GET("/menu", RequirePermission("gateway.menu.read"), wrapHandler(service.GetAllMenu))
	route.PUT("/menu", RequirePermission("gateway.menu.update"), wrapHandler(service.UpdateMenu))
	route.POST("/menu/create", RequirePermission("gateway.menu.create"), wrapHandler(service.CreateNewMenu))
	route.DELETE("/menu/:id", RequirePermission("gateway.menu.delete"), wrapHandler(service.DeleteMenu))
}
