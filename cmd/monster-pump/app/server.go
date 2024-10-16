package app

import (
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	genericapiserver "k8s.io/apiserver/pkg/server"
	logsapi "k8s.io/component-base/logs/api/v1"

	"github.com/rosas99/streaming/cmd/streaming-pump/app/options"
	"github.com/rosas99/streaming/internal/pkg/feature"
	"github.com/rosas99/streaming/internal/pump"
	"github.com/rosas99/streaming/pkg/app"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(feature.DefaultMutableFeatureGate))
}

const commandDesc = `Pump is a pluggable analytics purger to move Analytics generated by your streaming nodes to any back-end.

Find more streaming-pump information at:
    https://github.com/rosas99/streaming/blob/master/docs/guide/en-US/cmd/streaming-pump.md`

// NewApp creates an App object with default parameters.
func NewApp(name string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(name, "Launch a streaming pump server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithHealthCheckFunc(func() error {
			go opts.HealthOptions.ServeHealthCheck()
			return nil
		}),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func() error {
		cfg, err := opts.Config()
		if err != nil {
			return err
		}

		return Run(cfg, genericapiserver.SetupSignalHandler())
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(c *pump.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New()
	if err != nil {
		return err
	}

	return server.PrepareRun().Run(stopCh)
}
