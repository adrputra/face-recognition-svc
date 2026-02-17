package router

import "github.com/gin-gonic/gin"

func InitPermissionRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.permission

	route.GET("", RequirePermission("gateway.permission.read"), wrapHandler(service.GetAllPermissions))
	route.POST("", RequirePermission("gateway.permission.create"), wrapHandler(service.CreatePermission))
	route.PUT("", RequirePermission("gateway.permission.update"), wrapHandler(service.UpdatePermission))

	route.POST("/assign", RequirePermission("gateway.permission.assign"), wrapHandler(service.AssignRolePermissions))
	route.GET("/role/:id", RequirePermission("gateway.permission.read"), wrapHandler(service.GetRolePermissions))
}
