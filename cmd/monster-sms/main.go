package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/rosas99/streaming/cmd/streaming-sms/app"
)

func main() {
	app.NewApp("streaming-sms").Run()
}
