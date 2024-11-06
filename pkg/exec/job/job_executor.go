package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const defaultOnStopTimeout = 5 * time.Second

type Job struct {
	Config        ExecuteConfig
	Fn            func(ctx context.Context) error
	CloseFn       func(context.Context) error
	Name          string
	OnStopTimeout time.Duration
}

type Runner struct {
	logger *slog.Logger

	queue []Job
}

func NewRunner(logger *slog.Logger) *Runner {
	return &Runner{
		logger: logger,
		queue:  []Job{},
	}
}

func (r *Runner) AddJob(job Job) {
	r.queue = append(r.queue, job)
}

var ErrInvalidJobConfig = errors.New("invalid job config")

func (r *Runner) Run(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, job := range r.queue {
		switch cfg := job.Config.(type) {
		case IntervalConfig:
			if err := r.runIntervalFn(egCtx, eg, job, cfg); err != nil {
				return errors.Wrap(err, "failed to run interval job")
			}
		case ProcessConfig:
			r.runProcessFn(egCtx, eg, job)
		default:
			return ErrInvalidJobConfig
		}
	}

	return eg.Wait()
}

func (r *Runner) execStopFn(ctx context.Context, job Job) {
	if job.CloseFn != nil {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), max(job.OnStopTimeout, defaultOnStopTimeout))
		defer cancel()

		if err := job.CloseFn(ctx); err != nil {
			r.logger.Error("failed to execute on stop function",
				slog.String("job", job.Name),
				slog.String("error", err.Error()),
			)
		}
	}
}

func wrapNotCtxErr(ctx context.Context, err error, text string, args ...any) error {
	if errors.Is(err, ctx.Err()) {
		return nil
	}

	return errors.Wrapf(err, text, args...)
}
