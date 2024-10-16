package sms

import (
	genericoptions "github.com/rosas99/streaming/pkg/options"
)

func scheme(opts *genericoptions.TLSOptions) string {
	scheme := "http"
	if opts != nil && opts.UseTLS {
		scheme = "https"
	}

	return scheme
}
