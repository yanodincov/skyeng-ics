package executor

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (e *JobExecutor) runProcessFn(
	ctx context.Context,
	eg *errgroup.Group,
	job Job,
) {
	eg.Go(func() error {
		defer e.execStopFn(ctx, job)

		if err := job.Fn(ctx); err != nil {
			return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
		}

		return nil
	})
}
