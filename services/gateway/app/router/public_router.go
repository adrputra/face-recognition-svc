package router

import (
	model "github.com/adrputra/face-recognition-svc/gateway/app/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitPublicRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.user

	// PUBLIC
	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "pong",
			Data:    nil,
		})
	})

	// PUBLIC
	route.POST("/register", service.CreateNewUser)
	// PUBLIC
	route.POST("/login", service.Login)
}
