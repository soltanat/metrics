package client

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/soltanat/metrics/internal/logger"
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
	req.Body = io.NopCloser(&buf)
	req.ContentLength = int64(buf.Len())
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

type SignatureTransport struct {
	Transport http.RoundTripper
	Key       string
}

func (t *SignatureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body := req.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	teeReader := io.TeeReader(body, buf)

	hash := sha256.New()
	_, err := io.Copy(hash, teeReader)
	if err != nil {
		return nil, err
	}

	hash.Write([]byte(t.Key))
	req.Header.Set("HashSHA256", fmt.Sprintf("%x", hash.Sum(nil)))

	req.Body = io.NopCloser(buf)

	return t.Transport.RoundTrip(req)
}
