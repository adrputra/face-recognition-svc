package router

import (
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitPublicRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.user

	// PUBLIC
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "pong",
			Data:    nil,
		})
	})

	// PUBLIC
	route.POST("/register", wrapHandler(service.CreateNewUser))
	// PUBLIC
	route.POST("/login", wrapHandler(service.Login))
}
