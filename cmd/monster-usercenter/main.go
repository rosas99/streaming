package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/rosas99/streaming/cmd/streaming-usercenter/app"
)

func main() {
	app.NewApp("streaming-usercenter").Run()
}
