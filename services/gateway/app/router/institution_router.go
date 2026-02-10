package router

import "github.com/labstack/echo/v4"

func InitInstitutionRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.institution

	route.GET("", service.GetAllInstitution, RequirePermission("gateway.institution.read"))
	route.GET("/:id", service.GetInstitutionByID, RequirePermission("gateway.institution.read"))
	route.POST("", service.CreateNewInstitution, RequirePermission("gateway.institution.create"))
	route.PUT("", service.UpdateInstitution, RequirePermission("gateway.institution.update"))
	route.DELETE("/:id", service.DeleteInstitution, RequirePermission("gateway.institution.delete"))
}
