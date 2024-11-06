package httphelper

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

const (
	contentEncodingHeader  = "Content-Encoding"
	contentEncodingGzip    = "gzip"
	contentEncodingDeflate = "deflate"
	contentEncodingBr      = "br"
	contentEncodingZstd    = "zstd"
)

// DecodeHTTPResponseBody decodes HTTP response bytesBody based on Content-Encoding header.
func DecodeHTTPResponseBody(res *http.Response) ([]byte, error) {
	var (
		reader io.ReadCloser
		err    error
	)

	// Select the appropriate reader based on Content-Encoding
	switch res.Header.Get(contentEncodingHeader) {
	case contentEncodingGzip:
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create gzip reader")
		}
		defer reader.Close()

	case contentEncodingDeflate:
		// No explicit Close method needed for flate.NewReader, so wrap it in io.NopCloser
		reader = io.NopCloser(flate.NewReader(res.Body))
		defer reader.Close()

	case contentEncodingBr:
		reader = io.NopCloser(brotli.NewReader(res.Body))
		defer reader.Close()

	case contentEncodingZstd:
		zstdReader, err := zstd.NewReader(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create zstd reader")
		}

		defer zstdReader.Close()

		reader = zstdReader.IOReadCloser()

	default:
		// If no encoding, read the bytesBody as is
		reader = res.Body
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return data, nil
}
