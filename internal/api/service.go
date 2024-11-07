package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/yanodincov/skyeng-ics/pkg/mem"

	"github.com/yanodincov/skyeng-ics/pkg/exec/job"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar"
	"github.com/yanodincov/skyeng-ics/pkg/exec"
)

const port = 8080

type Service struct {
	cfg             *config.Config
	calendarService *calendar.Service
	log             *slog.Logger
	runner          *job.Runner
}

func NewService(
	cfg *config.Config,
	calendarService *calendar.Service,
	log *slog.Logger,
	runner *job.Runner,
) *Service {
	service := &Service{
		cfg:             cfg,
		calendarService: calendarService,
		log:             log,
		runner:          runner,
	}
	service.onStart()

	return service
}

func (s *Service) onStart() {
	server := s.createServer()

	s.runner.AddJob(job.Job{ //nolint:exhaustruct
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
		Config: job.ProcessConfig{},
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

	serializedCalendar := []byte(icsCalendar.Serialize())
	ctx.Success("text/calendar; charset=utf-8", serializedCalendar)
	s.log.Info("successfully sent calendar",
		slog.String("size", mem.GetHumanReadableSize(len(serializedCalendar))),
	)
}
