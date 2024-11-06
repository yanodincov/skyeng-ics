package config

import (
	"time"

	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

const minRefreshInterval = time.Second * 30

type Config struct {
	Skyeng SkyengConfig
	Worker WorkerConfig
	API    APIConfig
}

type SkyengConfig struct {
	Login    string `envconfig:"SKYENG_LOGIN"`
	Password string `envconfig:"SKYENG_PASSWORD"`
}

type WorkerConfig struct {
	RefreshInterval time.Duration `envconfig:"REFRESH_INTERVAL" default:"5m"`
}

type APIConfig struct {
	Port        int    `envconfig:"API_PORT" default:"8080"`
	RouteSuffix string `envconfig:"API_ROUTE_SUFFIX" `
}

func ParseConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	cfg.Worker.RefreshInterval = max(cfg.Worker.RefreshInterval, minRefreshInterval)

	return &cfg, nil
}