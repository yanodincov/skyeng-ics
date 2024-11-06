package skyeng

import (
	"net/http"
	"time"

	"github.com/yanodincov/skyeng-ics/internal/config"
)

const (
	maxIdleConnections = 10
	timeout            = 10 * time.Second
)

type Repository struct {
	cfg    *config.SkyengConfig
	client *http.Client
}

func NewRepository(
	cfg *config.SkyengConfig,
) *Repository {
	return &Repository{
		cfg: cfg,
		client: &http.Client{ //nolint:exhaustruct
			Timeout: timeout,
			Transport: &http.Transport{ //nolint:exhaustruct
				MaxIdleConns: maxIdleConnections,
			},
		},
	}
}
