package calendar

import (
	"context"
	"sync"
	"time"

	"github.com/yanodincov/skyeng-ics/pkg/executor"

	ics "github.com/arran4/golang-ical"
	"github.com/pkg/errors"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng"
	"github.com/yanodincov/skyeng-ics/internal/service/auth"
	"github.com/yanodincov/skyeng-ics/internal/service/calendar/factory"
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
	jobExecutor     *executor.JobExecutor

	calendar *ics.Calendar
	mx       sync.RWMutex
}

func NewService(
	cfg *config.Config,
	repository *skyeng.Repository,
	authService *auth.Service,
	calendarFactory *factory.Factory,
	jobExecutor *executor.JobExecutor,
) *Service {
	service := &Service{ //nolint:exhaustruct
		cfg:             cfg,
		repository:      repository,
		authService:     authService,
		calendarFactory: calendarFactory,
		jobExecutor:     jobExecutor,
	}
	service.onStart()

	return service
}

func (s *Service) GetCalendar(_ context.Context) (*ics.Calendar, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.calendar, nil
}

func (s *Service) onStart() {
	s.jobExecutor.AddJob(executor.Job{ //nolint:exhaustruct
		Name: "refresh calendar",
		Fn: func(ctx context.Context) error {
			return s.refreshCalendar(ctx)
		},
		Config: executor.IntervalConfig{
			Interval: s.cfg.Worker.RefreshInterval,
			Timeout:  refreshTimeout,
			Retries:  refreshRetries,
		},
	})
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
