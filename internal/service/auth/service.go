package auth

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yanodincov/skyeng-ics/internal/config"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng/meta"
	"github.com/yanodincov/skyeng-ics/pkg/worker"
)

const (
	authCookiesInterval = 24 * time.Hour
	authTimeout         = 30 * time.Second
	authRetry           = 5
)

type Service struct {
	cfg          *config.Config
	repository   *skyeng.Repository
	metaProvider *meta.Provider
	spec         skyeng.GetScheduleSpec
	mx           sync.RWMutex
}

func NewService(
	cfg *config.Config,
	repository *skyeng.Repository,
	metaProvider *meta.Provider,
) *Service {
	return &Service{ //nolint:exhaustruct
		cfg:          cfg,
		repository:   repository,
		metaProvider: metaProvider,
	}
}

func (s *Service) GetAuthorizedGetScheduleSpec(_ context.Context) skyeng.GetScheduleSpec {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.spec
}

func (s *Service) Run(sd worker.IShutdowner) error {
	return worker.RunWorker(sd, authCookiesInterval, //nolint:wrapcheck
		func(ctx context.Context) error {
			spec, err := s.generateGetScheduleSpec(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to generate get schedule spec")
			}

			s.mx.Lock()
			s.spec = *spec
			s.mx.Unlock()

			return nil
		},
		worker.WithExecRetries(authRetry),
		worker.WithExecTimeout(authTimeout),
		worker.WithAsyncErrWrap("failed to run auth service"),
	)
}

func (s *Service) generateGetScheduleSpec(ctx context.Context) (*skyeng.GetScheduleSpec, error) {
	ctx, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()

	reqMeta := s.metaProvider.GenerateSkyengMeta()

	csrfTokenRes, err := s.repository.GetCsrfToken(ctx, skyeng.GetCsrfTokenSpec{
		Headers: reqMeta.GetScheduleHeaders,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CSRF token")
	}

	loginRes, err := s.repository.Login(ctx, skyeng.LoginSpec{
		CsrfToken: csrfTokenRes.CsrfToken,
		Login:     s.cfg.Skyeng.Login,
		Password:  s.cfg.Skyeng.Password,
		Headers:   reqMeta.LoginHeaders,
		Cookies:   csrfTokenRes.Cookies,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	authCookiesRes, err := s.repository.GetAuthCookies(ctx, skyeng.GetAuthCookiesSpec{
		Cookies: loginRes.Cookies,
		Headers: reqMeta.GetScheduleHeaders,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth cookies")
	}

	userIDRes, err := s.repository.GetUserID(ctx, skyeng.GetUserIDSpec{
		Cookies: authCookiesRes.Cookies,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user ID")
	}

	return &skyeng.GetScheduleSpec{
		Cookies: userIDRes.Cookies,
		Headers: reqMeta.GetScheduleHeaders,
		UserID:  userIDRes.UserID,
	}, nil
}