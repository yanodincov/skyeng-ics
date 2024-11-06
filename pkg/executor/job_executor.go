package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Job struct {
	Config        ExecuteConfig
	Fn            func(ctx context.Context) error
	CloseFn       func(context.Context) error
	Name          string
	OnStopTimeout time.Duration
}

type JobExecutor struct {
	logger *slog.Logger

	jobQueue []Job
}

func NewJobExecutor(logger *slog.Logger) *JobExecutor {
	return &JobExecutor{
		logger:   logger,
		jobQueue: []Job{},
	}
}

func (e *JobExecutor) AddJob(job Job) {
	e.jobQueue = append(e.jobQueue, job)
}

var ErrInvalidJobConfig = errors.New("invalid job config")

func (e *JobExecutor) Run(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	for _, job := range e.jobQueue {
		switch cfg := job.Config.(type) {
		case IntervalConfig:
			if err := e.runIntervalFn(egCtx, eg, job, cfg); err != nil {
				return errors.Wrap(err, "failed to run interval job")
			}
		case ProcessConfig:
			e.runProcessFn(egCtx, eg, job)
		default:
			return ErrInvalidJobConfig
		}
	}

	return eg.Wait()
}

func (e *JobExecutor) execStopFn(ctx context.Context, job Job) {
	if job.CloseFn != nil {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), job.OnStopTimeout)
		defer cancel()

		if err := job.CloseFn(ctx); err != nil {
			e.logger.Error("failed to execute on stop function",
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
