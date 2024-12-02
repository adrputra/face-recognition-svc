package router

import "github.com/labstack/echo/v4"

func InitRoleRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.role

	route.GET("", service.GetAllRole)

	route.GET("/mapping", service.GetAllRoleMapping)
	route.POST("/create", service.CreateNewRole)
	route.POST("/mapping/create", service.CreateNewRoleMapping)

	route.GET("/menu", service.GetAllMenu)
	route.PUT("/menu", service.UpdateMenu)
	route.POST("/menu/create", service.CreateNewMenu)
	route.DELETE("/menu/:id", service.DeleteMenu)
}
