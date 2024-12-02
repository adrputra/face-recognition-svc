package main

import (
	"face-recognition-svc/app"
	"os"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	app.Start()
}