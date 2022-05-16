package exponent_retry

import (
	"time"
)

var DefaultOptions = &Options{
	// Default is 3 retries
	MaxRetries: DefaultMaxRetries,
	// Default is 8 milliseconds
	MinRetryBackoff: DefaultMinRetryBackoff,
	// Default is 512 milliseconds
	MaxRetryBackoff: DefaultMaxRetryBackoff,
}

type Options struct {
	// Maximum number of retries before giving up.
	// Default is 3 retries.
	MaxRetries int

	// Minimum backoff between each retry.
	// Default is 8 milliseconds.
	MinRetryBackoff time.Duration

	// Maximum backoff between each retry.
	// Default is 512 milliseconds.
	MaxRetryBackoff time.Duration

	// Whether to retry according to the error returned
	// Default is nil.
	Fn Fn
}

func (opt *Options) init() {
	if opt.MaxRetries <= 0 {
		opt.MaxRetries = 3
	}

	if opt.MinRetryBackoff <= 0 {
		opt.MinRetryBackoff = 8 * time.Millisecond
	}

	if opt.MaxRetryBackoff <= 0 {
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}
}
