package skyeng

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	httphelper "github.com/yanodincov/skyeng-ics/pkg/http-helper"
)

const loginURL = "https://id.skyeng.ru/frame/login-submit"

type LoginSpec struct {
	CsrfToken string
	Login     string
	Password  string
	Headers   map[string]string
	Cookies   []*http.Cookie
}

type LoginData struct {
	Cookies []*http.Cookie
}

func (r *Repository) Login(ctx context.Context, spec LoginSpec) (*LoginData, error) {
	formData := url.Values{}
	formData.Set("csrfToken", spec.CsrfToken)
	formData.Set("redirect", "https://skyeng.ru/")
	formData.Set("username", spec.Login)
	formData.Set("password", spec.Password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	for header, headerVal := range spec.Headers {
		req.Header.Set(header, headerVal)
	}

	for _, cookie := range spec.Cookies {
		req.AddCookie(cookie)
	}

	httpRes, err := r.client.Do(req)
	defer httpRes.Body.Close() //nolint:govet,staticcheck

	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", httpRes.StatusCode)
	}

	return &LoginData{
		Cookies: httphelper.MergeCookies(httpRes, spec.Cookies),
	}, nil
}
