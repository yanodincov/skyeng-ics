package job

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

func (r *Runner) runProcessFn(
	ctx context.Context,
	eg *errgroup.Group,
	job Job,
) {
	eg.Go(func() error {
		defer r.execStopFn(ctx, job)

		r.logger.Info("running process job", slog.String("job", job.Name))

		if err := job.Fn(ctx); err != nil {
			return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
		}

		return nil
	})
}
