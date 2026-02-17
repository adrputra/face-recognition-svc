package router

import (
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"
	"github.com/gin-gonic/gin"
)

type handlerWithError func(*gin.Context) error

func wrapHandler(handler handlerWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil && !c.Writer.Written() {
			_ = utils.LogError(c, err, nil)
		}
	}
}
