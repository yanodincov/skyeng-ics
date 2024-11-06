package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/yanodincov/skyeng-ics/internal/api"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng/meta"
	"github.com/yanodincov/skyeng-ics/internal/service/auth"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar/factory"
	"github.com/yanodincov/skyeng-ics/pkg/worker"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ //nolint:exhaustruct
		AddSource: false,
		Level:     slog.LevelInfo,
	}))

	if err := run(ctx, logger); err != nil {
		logger.Error("failed to run the application", slog.String("error", err.Error()))
	}
}

func run(ctx context.Context, log *slog.Logger) error {
	cfg, err := config.ParseConfig()
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	shutdowner := worker.NewShutdowner(ctx, log)
	defer shutdowner.Stop()

	metaProvider := meta.NewProvider()
	skyengRepository := skyeng.NewRepository(&cfg.Skyeng)

	authService := auth.NewService(cfg, skyengRepository, metaProvider)
	if err := authService.Run(shutdowner); err != nil {
		return errors.Wrap(err, "failed to run auth service")
	}

	calendarFactory := factory.NewFactory()

	calendarService := calendar.NewService(
		cfg,
		skyengRepository,
		authService,
		calendarFactory,
	)
	if err = calendarService.Run(shutdowner); err != nil {
		return errors.Wrap(err, "failed to run calendar service")
	}

	apiService := api.NewService(cfg, calendarService, log)
	if err = apiService.Run(shutdowner); err != nil {
		return errors.Wrap(err, "failed to run api service")
	}

	<-ctx.Done()

	return nil
}
