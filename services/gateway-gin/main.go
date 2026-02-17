package main

import (
	"os"

	"github.com/adrputra/face-recognition-svc/gateway-gin/app"
	"github.com/gin-gonic/gin"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	gin.SetMode(gin.ReleaseMode)
	// os.Setenv("ENV", "development")
	app.Start()
}
