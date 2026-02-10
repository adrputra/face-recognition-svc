package router

import "github.com/labstack/echo/v4"

func InitRoleRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.role

	route.GET("", service.GetAllRole, RequirePermission("gateway.role.read"))

	route.GET("/mapping", service.GetAllRoleMapping, RequirePermission("gateway.role_mapping.read"))
	route.POST("/create", service.CreateNewRole, RequirePermission("gateway.role.create"))
	route.POST("/mapping/create", service.CreateNewRoleMapping, RequirePermission("gateway.role_mapping.create"))
	route.PUT("/mapping", service.UpdateRoleMapping, RequirePermission("gateway.role_mapping.update"))
	route.DELETE("/mapping/:id", service.DeleteRoleMapping, RequirePermission("gateway.role_mapping.delete"))

	route.GET("/menu", service.GetAllMenu, RequirePermission("gateway.menu.read"))
	route.PUT("/menu", service.UpdateMenu, RequirePermission("gateway.menu.update"))
	route.POST("/menu/create", service.CreateNewMenu, RequirePermission("gateway.menu.create"))
	route.DELETE("/menu/:id", service.DeleteMenu, RequirePermission("gateway.menu.delete"))
}
