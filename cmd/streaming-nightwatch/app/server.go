package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/rosas99/streaming/cmd/streaming-nightwatch/app/options"
	"github.com/rosas99/streaming/internal/nightwatch"
	"github.com/rosas99/streaming/pkg/app"
)

const commandDesc = `The sms server is a standard, specification-compliant demo 
example of the streaming service.

Find more streaming-sms information at:
    https://"github.com/rosas99/streaming/blob/master/docs/guide/en-US/cmd/streaming-nightwatch.md`

// NewApp creates an App object with default parameters.
func NewApp(name string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp(name, "Launch a streaming nightwatch server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
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
func Run(c *nightwatch.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New()
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}
