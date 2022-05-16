package exponent_retry

import "time"

const (
	DefaultMaxRetries      = 3
	DefaultMinRetryBackoff = 8 * time.Millisecond
	DefaultMaxRetryBackoff = 512 * time.Millisecond
)
