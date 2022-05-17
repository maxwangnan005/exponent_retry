package exponent_retry

import (
	"context"
	"io"
	"time"
)

func retryBackoff(retry int, minBackoff, maxBackoff time.Duration) time.Duration {
	if retry < 0 {
		panic("not reached")
	}
	// not reached
	if minBackoff == 0 {
		return 0
	}

	d := minBackoff << uint(retry-1)
	if d < minBackoff {
		return maxBackoff
	}

	if d > maxBackoff || d < minBackoff {
		d = maxBackoff
	}

	return d
}

func sleep(ctx context.Context, dur time.Duration) error {
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type Fn func(err error) bool

func shouldRetry(err error, fn Fn) bool {
	switch err {
	case io.EOF, io.ErrUnexpectedEOF:
		return true
	case nil, context.Canceled, context.DeadlineExceeded:
		return false
	}

	if fn != nil {
		return fn(err)
	}

	return true
}
