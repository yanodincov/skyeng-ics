package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

const minRefreshInterval = time.Second * 30

type Config struct {
	Skyeng SkyengConfig
	API    APIConfig
	Worker WorkerConfig
}

type SkyengConfig struct {
	Login    string `envconfig:"SKYENG_LOGIN"`
	Password string `envconfig:"SKYENG_PASSWORD"`
}

type WorkerConfig struct {
	RefreshInterval time.Duration `envconfig:"REFRESH_INTERVAL" default:"5m"`
}

type APIConfig struct {
	RouteSuffix string `envconfig:"API_ROUTE_SUFFIX"`
}

var ErrEmptySkyengCredentials = errors.New("empty skyeng credentials")

func ParseConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	cfg.Worker.RefreshInterval = max(cfg.Worker.RefreshInterval, minRefreshInterval)
	if cfg.Skyeng.Login == "" || cfg.Skyeng.Password == "" {
		return nil, ErrEmptySkyengCredentials
	}

	return &cfg, nil
}
