// Package options contains flags and options for initializing an apiserver
package options

import (
	"github.com/rosas99/streaming/internal/nightwatch"
	"github.com/rosas99/streaming/internal/pkg/feature"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/rosas99/streaming/pkg/app"
	"github.com/rosas99/streaming/pkg/log"
	genericoptions "github.com/rosas99/streaming/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting streaming-nightwatch server.
	UserAgent = "streaming-nightwatch"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains state for master/api server.
type Options struct {
	GRPCOptions  *genericoptions.GRPCOptions  `json:"grpc" mapstructure:"grpc"`
	HTTPOptions  *genericoptions.HTTPOptions  `json:"http" mapstructure:"http"`
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	//Redis options for configuring Redis related options.
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	Log          *log.Options                 `json:"log" mapstructure:"log"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
		RedisOptions: genericoptions.NewRedisOptions(),
		Log:          log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	//Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	//arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

// Complete completes all the required options.
func (o *Options) Complete() error {

	//_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	errs := []error{}

	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up streaming-fakeserver config with options.
func (o *Options) ApplyTo(c *nightwatch.Config) error {
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	return nil
}

// Config return a streaming-fakeserver config object.
func (o *Options) Config() (*nightwatch.Config, error) {
	c := &nightwatch.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
