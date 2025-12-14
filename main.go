package main

import (
	"face-recognition-svc/app"
	"os"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	os.Setenv("ENV", "development")
	app.Start()
}
