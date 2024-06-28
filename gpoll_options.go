package gpoll

import "gpoll/conn"

// Options params
type Options struct {
	OnRequest conn.ConnectionHandler
}
type Option func(options *Options)

func WithOnRequest(handler conn.ConnectionHandler) Option {
	return func(options *Options) {
		options.OnRequest = handler
	}
}

func completeOptions(setters ...Option) *Options {
	o := &Options{}
	for _, setter := range setters {
		setter(o)
	}
	return o
}
