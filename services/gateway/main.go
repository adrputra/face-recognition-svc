package main

import (
	"face-recognition-svc/gateway/app"
	"os"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	// os.Setenv("ENV", "development")
	app.Start()
}
