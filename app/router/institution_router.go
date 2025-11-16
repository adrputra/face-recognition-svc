package router

import "github.com/labstack/echo/v4"

func InitInstitutionRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.institution

	route.GET("", service.GetAllInstitution)
	route.GET("/:id", service.GetInstitutionByID)
	route.POST("", service.CreateNewInstitution)
	route.PUT("", service.UpdateInstitution)
	route.DELETE("/:id", service.DeleteInstitution)
}
