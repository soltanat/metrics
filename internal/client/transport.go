package client

import (
	"bytes"
	"compress/gzip"
	"github.com/soltanat/metrics/internal/logger"
	"io"
	"net/http"
	"time"
)

type GzipTransport struct {
	Transport http.RoundTripper
}

func (t *GzipTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := io.Copy(gw, req.Body)
	if err != nil {
		return nil, err
	}
	err = gw.Close()
	if err != nil {
		return nil, err
	}

	err = req.Body.Close()
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(&buf)        // Set the compressed body to the request
	req.ContentLength = int64(buf.Len()) // Update the Content-Length header
	req.Header.Set("Content-Encoding", "gzip")

	return t.Transport.RoundTrip(req)
}

type LoggingTransport struct {
	Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.Transport.RoundTrip(req)
	latency := time.Since(start)

	l := logger.Get()

	l.Info().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Dur("latency", latency).
		Msg("Request")
	return resp, err
}
