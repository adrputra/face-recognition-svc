package router

import "github.com/labstack/echo/v4"

func InitPresenceRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.presence

	route.GET("", service.ListPresence, RequirePermission("gateway.presence.read"))
	route.GET("/:id", service.GetPresence, RequirePermission("gateway.presence.read"))
	route.POST("", service.CreatePresence, RequirePermission("gateway.presence.create"))
	route.PUT("/:id", service.UpdatePresence, RequirePermission("gateway.presence.update"))
	route.DELETE("/:id", service.DeletePresence, RequirePermission("gateway.presence.delete"))
}
