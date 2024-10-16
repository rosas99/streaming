// Package options contains flags and options for initializing an usercenter server
package options

import (
	"github.com/rosas99/streaming/internal/usercenter"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/rosas99/streaming/pkg/app"
	"github.com/rosas99/streaming/pkg/log"
	genericoptions "github.com/rosas99/streaming/pkg/options"
)

const (
	// UserAgent is the userAgent name when starting streaming-sms server.
	UserAgent = "streaming-usercenter"
)

var _ app.CliOptions = (*Options)(nil)

// Options contains state for master/api server.
type Options struct {
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	//TLSOptions   *genericoptions.TLSOptions     `json:"tls" mapstructure:"tls"`
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	//Redis options for configuring Redis related options.
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	Log          *log.Options                 `json:"log" mapstructure:"log"`
}

// NewOptions returns initialized Options.
func NewOptions() *Options {
	o := &Options{
		GRPCOptions: genericoptions.NewGRPCOptions(),
		HTTPOptions: genericoptions.NewHTTPOptions(),
		//TLSOptions:    genericoptions.NewTLSOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
		RedisOptions: genericoptions.NewRedisOptions(),
		Log:          log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	//o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	//o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	//o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	//o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	//o.Log.AddFlags(fss.FlagSet("log"))
	//o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	//fs := fss.FlagSet("misc")
	//feature.DefaultMutableFeatureGate.AddFlag(fs)

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
	//errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	//errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up streaming-fakeserver config with options.
func (o *Options) ApplyTo(c *usercenter.Config) error {
	c.GRPCOptions = o.GRPCOptions
	c.HTTPOptions = o.HTTPOptions
	//c.TLSOptions = o.TLSOptions
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	return nil
}

// Config return a streaming-usercenter server config object.
func (o *Options) Config() (*usercenter.Config, error) {
	c := &usercenter.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
