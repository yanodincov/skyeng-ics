package job

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (r *Runner) runProcessFn(
	ctx context.Context,
	eg *errgroup.Group,
	job Job,
) {
	eg.Go(func() error {
		defer r.execStopFn(ctx, job)

		if err := job.Fn(ctx); err != nil {
			return wrapNotCtxErr(ctx, err, "failed to execute worker function for job %s", job.Name)
		}

		return nil
	})
}
