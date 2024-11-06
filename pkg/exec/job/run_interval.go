package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"github.com/yanodincov/skyeng-ics/pkg/exec"

	"golang.org/x/sync/errgroup"
)

func (r *Runner) runIntervalFn(
	ctx context.Context,
	eg *errgroup.Group,
	job Job,
	cfg IntervalConfig,
) error {
	if err := r.executeIntervalJob(ctx, job, cfg); err != nil {
		r.execStopFn(ctx, job)

		return wrapNotCtxErr(ctx, err, "failed first execute worker job '%s'", job.Name)
	}

	r.logger.Info("running interval job", slog.String("job", job.Name))

	eg.Go(func() error {
		defer r.execStopFn(ctx, job)

		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				if err := r.executeIntervalJob(ctx, job, cfg); err != nil {
					return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
				}
			}
		}
	})

	return nil
}

func (r *Runner) executeIntervalJob(ctx context.Context, job Job, cfg IntervalConfig) error {
	start := time.Now()

	if err := exec.Retry(ctx, cfg.Retries, job.Fn, exec.WithExecuteTimeout(cfg.Timeout)); err != nil {
		r.logger.Error("failed to execute worker function",
			slog.String("job", job.Name),
			slog.String("error", err.Error()),
		)

		return errors.Wrap(err, "failed to execute worker function")
	}

	r.logger.Info("job executed",
		slog.String("job", job.Name),
		slog.Duration("duration", time.Since(start)),
	)

	return nil
}
