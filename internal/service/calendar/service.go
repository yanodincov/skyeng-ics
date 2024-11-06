package calendar

import (
	"context"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/pkg/errors"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng"
	"github.com/yanodincov/skyeng-ics/internal/service/auth"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar/factory"
	"github.com/yanodincov/skyeng-ics/pkg/worker"
)

const (
	refreshTimeout = 30 * time.Second
	refreshRetries = 5
)

type Service struct {
	cfg             *config.Config
	authService     *auth.Service
	repository      *skyeng.Repository
	calendarFactory *factory.Factory
	calendar        *ics.Calendar
	mx              sync.RWMutex
}

func NewService(
	cfg *config.Config,
	repository *skyeng.Repository,
	authService *auth.Service,
	calendarFactory *factory.Factory,
) *Service {
	return &Service{ //nolint:exhaustruct
		cfg:             cfg,
		repository:      repository,
		authService:     authService,
		calendarFactory: calendarFactory,
	}
}

func (s *Service) GetCalendar(_ context.Context) (*ics.Calendar, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.calendar, nil
}

func (s *Service) Run(sd worker.IShutdowner) error {
	return worker.RunWorker( //nolint:wrapcheck
		sd,
		s.cfg.Worker.RefreshInterval,
		s.refreshCalendar,
		worker.WithExecTimeout(refreshTimeout),
		worker.WithExecRetries(refreshRetries),
		worker.WithAsyncErrWrap("failed to run calendar service"),
	)
}

func (s *Service) refreshCalendar(ctx context.Context) error {
	getScheduleData, err := s.repository.GetSchedule(ctx, s.authService.GetAuthorizedGetScheduleSpec(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to get schedule")
	}

	calendar, err := s.calendarFactory.CreateCalendarFromLessons(ctx, getScheduleData.Lessons)
	if err != nil {
		return errors.Wrap(err, "create calendar")
	}

	s.mx.Lock()
	s.calendar = calendar
	s.mx.Unlock()

	return nil
}
