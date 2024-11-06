package worker

import (
	"context"
	"log/slog"
)

type IShutdowner interface {
	Terminate(err error)
	GetContext() context.Context
	Stop()
}

type shutdowner struct {
	logger *slog.Logger

	ctx       context.Context //nolint:containedctx
	ctxCancel context.CancelFunc
}

func NewShutdowner(ctx context.Context, logger *slog.Logger) IShutdowner {
	ctx, ctxCancel := context.WithCancel(ctx)

	return &shutdowner{
		logger:    logger,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
}

func (s *shutdowner) GetContext() context.Context {
	return s.ctx
}

func (s *shutdowner) Terminate(err error) {
	s.logger.Error("terminate app", slog.String("error", err.Error()))
	s.ctxCancel()
}

func (s *shutdowner) Stop() {
	s.ctxCancel()
}
