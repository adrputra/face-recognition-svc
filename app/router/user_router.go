package router

import "github.com/labstack/echo/v4"

func InitUserRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.user

	route.GET("", service.GetAllUser)
	route.GET("/detail/:id", service.GetUserDetail)

	route.GET("/institutions", service.GetInstitutionList)
}
