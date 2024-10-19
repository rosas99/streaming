package main

import (
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration

	"github.com/rosas99/streaming/cmd/streaming-pump/app"
)

func main() {
	app.NewApp("streaming-pump").Run()
}
