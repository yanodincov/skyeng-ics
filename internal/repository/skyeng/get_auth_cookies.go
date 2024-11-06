package skyeng

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"
)

const (
	getAuthCookiesURL = "https://id.skyeng.ru/user-api/v1/auth/jwt"
)

type GetAuthCookiesSpec struct {
	Headers map[string]string
	Cookies []*http.Cookie
}

type GetAuthCookiesData struct {
	Cookies []*http.Cookie
}

func (r *Repository) GetAuthCookies(ctx context.Context, spec GetAuthCookiesSpec) (*GetAuthCookiesData, error) {
	req, err := httphelper.NewRequest(ctx, http.MethodPost, getAuthCookiesURL,
		httphelper.WithCookies(spec.Cookies),
		httphelper.WithHeaders(spec.Headers),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	httpRes, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", httpRes.StatusCode)
	}

	return &GetAuthCookiesData{
		Cookies: httphelper.MergeCookies(httpRes, spec.Cookies),
	}, nil
}
