package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/rosas99/streaming/cmd/streaming-nightwatch/app"
)

func main() {
	app.NewApp("streaming-nightwatch").Run()
}
