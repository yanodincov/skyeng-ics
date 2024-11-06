package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar"
	"github.com/yanodincov/skyeng-ics/pkg/worker"
)

type Service struct {
	cfg             *config.Config
	calendarService *calendar.Service
	log             *slog.Logger
}

func NewService(
	cfg *config.Config,
	calendarService *calendar.Service,
	log *slog.Logger,
) *Service {
	return &Service{
		cfg:             cfg,
		calendarService: calendarService,
		log:             log,
	}
}

func (s *Service) Run(shutdowner worker.IShutdowner) error {
	route := "/calendar.ics"
	if s.cfg.API.RouteSuffix != "" {
		route = "/" + strings.TrimPrefix(strings.TrimSuffix(s.cfg.API.RouteSuffix, "/"), "/") + route
	}

	server := &fasthttp.Server{ //nolint:exhaustruct
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case route:
				s.getCalendar(ctx)
			default:
				ctx.Error("not found", http.StatusNotFound)
			}
		},
	}

	go func() {
		<-shutdowner.GetContext().Done()

		if err := server.Shutdown(); err != nil {
			s.log.Error("failed to shutdown server", slog.String("error", err.Error()))
		}

		s.log.Info("server stopped")
	}()

	s.log.Info("starting server", slog.Int("port", s.cfg.API.Port))

	if err := server.ListenAndServe(":" + strconv.Itoa(s.cfg.API.Port)); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "failed to start server")
	}

	return nil
}

func (s *Service) getCalendar(ctx *fasthttp.RequestCtx) {
	icsCalendar, err := s.calendarService.GetCalendar(ctx)
	if err != nil {
		ctx.Error("failed to get calendar", http.StatusInternalServerError)

		return
	}

	ctx.Success("text/calendar; charset=utf-8", []byte(icsCalendar.Serialize()))
}
