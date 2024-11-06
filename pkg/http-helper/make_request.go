package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type NewRequestOpt func(*newRequestOpts)

func WithBytesBody(body []byte) NewRequestOpt {
	return func(opts *newRequestOpts) {
		opts.bytesBody = body
	}
}

func WithHeaders(headers map[string]string) NewRequestOpt {
	return func(opts *newRequestOpts) {
		opts.headers = headers
	}
}

func WithCookies(cookies []*http.Cookie) NewRequestOpt {
	return func(opts *newRequestOpts) {
		opts.cookies = cookies
	}
}

func WithJSONBody(body any) NewRequestOpt {
	return func(opts *newRequestOpts) {
		opts.jsonBody = body
	}
}

type newRequestOpts struct {
	method    string
	url       string
	jsonBody  any
	bytesBody []byte
	headers   map[string]string
	cookies   []*http.Cookie
}

func NewRequest(
	ctx context.Context,
	method string,
	url string,
	opts ...NewRequestOpt,
) (*http.Request, error) {
	options := newRequestOpts{ //nolint:exhaustruct
		method: method,
		url:    url,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.jsonBody != nil {
		bodyBytes, err := json.Marshal(options.jsonBody)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}

		options.bytesBody = bodyBytes
	}

	var bodyReader io.Reader
	if len(options.bytesBody) > 0 {
		bodyReader = bytes.NewReader(options.bytesBody)
	}

	req, err := http.NewRequestWithContext(ctx, options.method, options.url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	for header, headerVal := range options.headers {
		req.Header.Set(header, headerVal)
	}

	for _, cookie := range options.cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}
