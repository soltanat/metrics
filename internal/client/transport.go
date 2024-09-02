package client

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/soltanat/metrics/internal/logger"
)

// GzipTransport
// Транспорт для http клиента с gzip сжатием тела запроса
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

// LoggingTransport
// Транспорт для http клиента с логированием
type LoggingTransport struct {
	Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.Transport.RoundTrip(req)
	latency := time.Since(start)

	l := logger.Get()

	var statusCode string
	if resp != nil {
		statusCode = resp.Status
	}

	l.Info().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Dur("latency", latency).
		Str("status", statusCode).
		Msg("Request")
	return resp, err
}

// SignatureTransport
// Транспорт для http клиента с подписью запроса
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

type RSAEncryptionTransport struct {
	Transport http.RoundTripper
	Key       *rsa.PublicKey
}

func NewRSAEncryptionTransport(transport http.RoundTripper, key []byte) (*RSAEncryptionTransport, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("invalid block type")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return &RSAEncryptionTransport{
		Transport: transport,
		Key:       pubKey,
	}, nil
}

func (t *RSAEncryptionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, t.Key, bodyBytes)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(cipherText))

	return t.Transport.RoundTrip(req)
}

type XRealIPTransport struct {
	Transport http.RoundTripper
}

func (t *XRealIPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//req.Header.Set("X-Real-IP", req.RemoteAddr)
	return t.Transport.RoundTrip(req)
}
