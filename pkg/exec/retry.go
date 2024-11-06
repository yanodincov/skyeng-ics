package exec

import (
	"context"
	"time"
)

type retryParams struct {
	executeTimeout time.Duration
	retryDelay     time.Duration
}

type OptionFn func(*retryParams)

func WithExecuteTimeout(timeout time.Duration) OptionFn {
	return func(p *retryParams) {
		p.executeTimeout = timeout
	}
}

func WithRetryDelay(delay time.Duration) OptionFn {
	return func(p *retryParams) {
		p.retryDelay = delay
	}
}

func Retry(ctx context.Context, times int, fn func(ctx context.Context) error, opts ...OptionFn) error {
	var params retryParams
	for _, opt := range opts {
		opt(&params)
	}

	var err error

	for range times {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err = execWithTimeout(ctx, fn, params.executeTimeout); err == nil {
			return nil
		}

		if params.retryDelay > 0 {
			select {
			case <-time.After(params.retryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return err
}

func execWithTimeout(ctx context.Context, fn func(ctx context.Context) error, timeout time.Duration) error {
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(execCtx)
}
