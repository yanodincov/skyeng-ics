package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/yanodincov/skyeng-ics/pkg/exec"

	"golang.org/x/sync/errgroup"
)

func (e *JobExecutor) runIntervalFn(
	ctx context.Context,
	eg *errgroup.Group,
	job Job,
	cfg IntervalConfig,
) error {
	if err := exec.Retry(ctx, cfg.Retries, job.Fn, exec.WithExecuteTimeout(cfg.Timeout)); err != nil {
		e.execStopFn(ctx, job)

		return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
	}

	eg.Go(func() error {
		defer e.execStopFn(ctx, job)

		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				start := time.Now()

				if err := exec.Retry(ctx, cfg.Retries, job.Fn, exec.WithExecuteTimeout(cfg.Timeout)); err != nil {
					e.logger.Error("failed to execute worker function",
						slog.String("job", job.Name),
						slog.String("error", err.Error()),
					)

					return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
				}

				e.logger.Info("job executed",
					slog.String("job", job.Name),
					slog.Duration("duration", time.Since(start)),
				)
			}
		}
	})

	return nil
}
