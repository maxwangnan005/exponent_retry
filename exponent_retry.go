package exponent_retry

import (
	"context"
)

var DefaultExponentRetry = &ExponentRetry{
	opt: DefaultOptions,
}

type Cmd func() error
type CmdWithReturn func() (interface{}, error)

func Do(ctx context.Context, cmd Cmd) error {
	return DefaultExponentRetry.Do(ctx, cmd)
}

func DoWithReturn(ctx context.Context, cmdWithReturn CmdWithReturn) (interface{}, error) {
	return DefaultExponentRetry.DoWithReturn(ctx, cmdWithReturn)
}

type ExponentRetry struct {
	opt *Options
}

func New(opts *Options) *ExponentRetry {
	opts.init()

	er := &ExponentRetry{
		opt: opts,
	}

	return er
}

func (er *ExponentRetry) Do(ctx context.Context, cmd Cmd) error {
	var lastErr error

	for attempt := 0; attempt <= er.opt.MaxRetries; attempt++ {
		currentAttempt := attempt

		retry, err := er.do(ctx, cmd, currentAttempt)
		if err == nil || !retry {
			return err
		}

		lastErr = err
	}

	return lastErr
}

func (er *ExponentRetry) do(ctx context.Context, cmd Cmd, attempt int) (bool, error) {
	if attempt > 0 { // Not the first time, this time is a retry
		if err := sleep(ctx, retryBackoff(attempt, er.opt.MinRetryBackoff, er.opt.MaxRetryBackoff)); err != nil {
			return false, err
		}
	}

	err := cmd()
	if err == nil {
		return false, nil
	}

	retry := shouldRetry(err, er.opt.Fn)

	return retry, err
}

func (er *ExponentRetry) DoWithReturn(ctx context.Context, cmdWithReturn CmdWithReturn) (interface{}, error) {
	var lastErr error
	var result interface{}

	for attempt := 0; attempt <= er.opt.MaxRetries; attempt++ {
		currentAttempt := attempt

		retry, res, err := er.doWithReturn(ctx, cmdWithReturn, currentAttempt)
		if err == nil || !retry {
			return res, err
		}

		lastErr = err
		result = res
	}

	return result, lastErr
}

func (er *ExponentRetry) doWithReturn(ctx context.Context, cmdWithReturn CmdWithReturn, attempt int) (bool, interface{}, error) {
	if attempt > 0 { // Not the first time, this time is a retry
		if err := sleep(ctx, retryBackoff(attempt, er.opt.MinRetryBackoff, er.opt.MaxRetryBackoff)); err != nil {
			return false, nil, err
		}
	}

	res, err := cmdWithReturn()
	if err == nil {
		return false, res, nil
	}

	retry := shouldRetry(err, er.opt.Fn)

	return retry, res, err
}
