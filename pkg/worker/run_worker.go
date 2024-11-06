package worker

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type runWorkerOpts struct {
	asyncErrWrap string
	execTimeout  time.Duration
	execRetries  int
}

type RunWorkerOption func(*runWorkerOpts)

func WithExecTimeout(timeout time.Duration) RunWorkerOption {
	return func(o *runWorkerOpts) {
		o.execTimeout = timeout
	}
}

func WithExecRetries(retries int) RunWorkerOption {
	return func(o *runWorkerOpts) {
		o.execRetries = retries
	}
}

func WithAsyncErrWrap(text string) RunWorkerOption {
	return func(o *runWorkerOpts) {
		o.asyncErrWrap = text
	}
}

func RunWorker(
	shutdowner IShutdowner,
	interval time.Duration,
	fn func(ctx context.Context) error,
	opts ...RunWorkerOption,
) error {
	opt := runWorkerOpts{} //nolint:exhaustruct
	for _, optFn := range opts {
		optFn(&opt)
	}

	opt.execRetries = max(opt.execRetries, 1)

	err := execWorkerFn(shutdowner.GetContext(), fn, opt.execTimeout, opt.execRetries)
	if err != nil {
		return errors.Wrap(err, "failed to execute worker function")
	}

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-shutdowner.GetContext().Done():
				return
			case <-ticker.C:
				err := execWorkerFn(shutdowner.GetContext(), fn, opt.execTimeout, opt.execRetries)
				if err != nil {
					if opt.asyncErrWrap != "" {
						err = errors.Wrap(err, opt.asyncErrWrap)
					}

					shutdowner.Terminate(errors.Wrap(err, "failed to execute worker function"))

					return
				}
			}
		}
	}()

	return nil
}

func execWorkerFn(
	ctx context.Context,
	fn func(ctx context.Context) error,
	timeout time.Duration,
	retries int,
) error {
	var err error

	for range retries {
		ctx, cancel := ctx, func() {}
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
		}

		err = fn(ctx)

		cancel()

		if err == nil {
			break
		}
	}

	return err
}
