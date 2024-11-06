package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/yanodincov/skyeng-ics/pkg/exec"
	"github.com/yanodincov/skyeng-ics/pkg/executor"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar"
)

const port = 8080

type Service struct {
	cfg             *config.Config
	calendarService *calendar.Service
	log             *slog.Logger
	jobExecutor     *executor.JobExecutor
}

func NewService(
	cfg *config.Config,
	calendarService *calendar.Service,
	log *slog.Logger,
	jobExecutor *executor.JobExecutor,
) *Service {
	service := &Service{
		cfg:             cfg,
		calendarService: calendarService,
		log:             log,
		jobExecutor:     jobExecutor,
	}
	service.onStart()

	return service
}

func (s *Service) onStart() {
	server := s.createServer()

	s.jobExecutor.AddJob(executor.Job{ //nolint:exhaustruct
		Name: "start api server",
		Fn: func(_ context.Context) error {
			if err := server.ListenAndServe(":" + strconv.Itoa(port)); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				return errors.Wrap(err, "failed to listen http server")
			}

			return nil
		},
		CloseFn: func(ctx context.Context) error {
			return exec.WithCtx(ctx, func() error {
				return errors.Wrap(server.Shutdown(), "failed to shutdown server")
			})
		},
		Config: executor.ProcessConfig{},
	})

	s.log.Info("starting server",
		slog.Int("port", port),
		slog.String("route", s.getRoute()),
	)
}

func (s *Service) createServer() *fasthttp.Server {
	route := s.getRoute()

	return &fasthttp.Server{ //nolint:exhaustruct
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case route:
				s.handlerGetCalendar(ctx)
			default:
				ctx.Error("not found", http.StatusNotFound)
			}
		},
	}
}

func (s *Service) getRoute() string {
	route := "/calendar.ics"
	if s.cfg.API.RouteSuffix != "" {
		route = "/" + strings.TrimPrefix(strings.TrimSuffix(s.cfg.API.RouteSuffix, "/"), "/") + route
	}

	return route
}

func (s *Service) handlerGetCalendar(ctx *fasthttp.RequestCtx) {
	icsCalendar, err := s.calendarService.GetCalendar(ctx)
	if err != nil {
		ctx.Error("failed to get calendar", http.StatusInternalServerError)

		return
	}

	ctx.Success("text/calendar; charset=utf-8", []byte(icsCalendar.Serialize()))
}
