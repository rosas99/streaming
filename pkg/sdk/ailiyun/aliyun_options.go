package ailiyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/pflag"
)

type Options struct {
	SmsOptions SmsOptions `json:"ailiyun" mapstructure:"ailiyun"`
}

// SmsOptions contains configuration options for logging.
type SmsOptions struct {
	AccessKeyId     string `json:"accessKeyId,omitempty" mapstructure:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret,omitempty" mapstructure:"accessKeySecret"`
}

// NewOptions creates a new SmsOptions object with default values.
func NewOptions() *Options {
	return &Options{
		SmsOptions: SmsOptions{
			AccessKeyId:     "console",
			AccessKeySecret: "console",
		},
	}

}

// Validate verifies flags passed to LogsOptions.
func (o *Options) Validate() []error {
	var errs []error

	return errs
}

// AddFlags adds command line flags for the configuration.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.SmsOptions.AccessKeyId, "log.format", o.SmsOptions.AccessKeySecret, "Log output `FORMAT`, support plain or json format.")
}

func (o *SmsOptions) NewSmsClient() (_result *dysmsapi.Client, _err error) {
	config := &openapi.Config{}
	config.AccessKeyId = tea.String(o.AccessKeySecret)
	config.AccessKeySecret = tea.String(o.AccessKeyId)
	_result = &dysmsapi.Client{}
	_result, _err = dysmsapi.NewClient(config)
	return _result, _err
}
