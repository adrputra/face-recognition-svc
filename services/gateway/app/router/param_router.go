package router

import "github.com/labstack/echo/v4"

func InitParamRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.param

	route.GET("/:id", service.GetParameterByKey, RequirePermission("gateway.param.read"))
	route.GET("", service.GetAllParam, RequirePermission("gateway.param.read"))
	route.POST("", service.InsertNewParam, RequirePermission("gateway.param.create"))
	route.PUT("", service.UpdateParam, RequirePermission("gateway.param.update"))
	route.DELETE("/:id", service.DeleteParam, RequirePermission("gateway.param.delete"))
}
