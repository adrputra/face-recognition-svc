package router

import "github.com/labstack/echo/v4"

func InitParamRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.param

	route.GET("/:id", service.GetParameterByKey)
	route.GET("", service.GetAllParam)
	route.POST("", service.InsertNewParam)
	route.PUT("", service.UpdateParam)
	route.DELETE("/:id", service.DeleteParam)
}
