package router

import "github.com/labstack/echo/v4"

func InitPermissionRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.permission

	route.GET("", service.GetAllPermissions)
	route.POST("", service.CreatePermission)
	route.PUT("", service.UpdatePermission)

	route.POST("/assign", service.AssignRolePermissions)
	route.GET("/role/:id", service.GetRolePermissions)
}
