package skyeng

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"
)

const (
	getCsrfTokenURL      = "https://id.skyeng.ru/login?redirect=https%3A%2F%2Fskyeng.ru%2F" //nolint:gosec
	getCsfrInputSelector = `input[type="hidden"][name="csrfToken"]`
)

type GetCsrfTokenSpec struct {
	Headers map[string]string
}

type GetCsrfTokenData struct {
	CsrfToken string
	Cookies   []*http.Cookie
}

func (r *Repository) GetCsrfToken(ctx context.Context, spec GetCsrfTokenSpec) (*GetCsrfTokenData, error) {
	req, err := httphelper.NewRequest(ctx, http.MethodGet, getCsrfTokenURL,
		httphelper.WithHeaders(spec.Headers),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	httpRes, err := r.client.Do(req)
	defer httpRes.Body.Close() //nolint:govet, staticcheck

	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", httpRes.StatusCode)
	}

	decodedBody, err := httphelper.DecodeHTTPResponseBody(httpRes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	csrfToken, err := httphelper.GetHTMLElemValue(decodedBody, getCsfrInputSelector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CSRF token")
	}

	return &GetCsrfTokenData{
		CsrfToken: csrfToken,
		Cookies:   httpRes.Cookies(),
	}, nil
}
