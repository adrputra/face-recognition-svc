package router

import "github.com/labstack/echo/v4"

func InitPermissionRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.permission

	route.GET("", service.GetAllPermissions, RequirePermission("gateway.permission.read"))
	route.POST("", service.CreatePermission, RequirePermission("gateway.permission.create"))
	route.PUT("", service.UpdatePermission, RequirePermission("gateway.permission.update"))

	route.POST("/assign", service.AssignRolePermissions, RequirePermission("gateway.permission.assign"))
	route.GET("/role/:id", service.GetRolePermissions, RequirePermission("gateway.permission.read"))
}
