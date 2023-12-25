package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/soltanat/metrics/internal/handler"
	"io"
	"net/http"
	"net/url"

	"github.com/soltanat/metrics/internal/model"
)

const (
	gaugeEndpointPrefix   = "/update/gauge"
	counterEndpointPrefix = "/update/counter"
	updateEndpointPrefix  = "/update123"
)

var errValidationName = fmt.Errorf("min name len 1")

type errHTTP struct {
	Err error
}

func (e errHTTP) Error() string {
	return e.Err.Error()
}

type errUnexpectedResponse struct {
	StatusCode int
	Message    []byte
}

func (e errUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response: %d, %s", e.StatusCode, e.Message)
}

type Client struct {
	address string
}

func New(address string) *Client {
	return &Client{
		address,
	}
}

func (c *Client) Send(m *model.Metric) error {
	if m.Name == "" {
		return errValidationName
	}

	var reqURL string
	switch m.Type {
	case model.MetricTypeGauge:
		reqURL, _ = url.JoinPath(
			c.address, gaugeEndpointPrefix, m.Name, m.ValueAsString(),
		)
	case model.MetricTypeCounter:
		reqURL, _ = url.JoinPath(
			c.address, counterEndpointPrefix, m.Name, m.ValueAsString(),
		)
	}

	return c.makeRequest(reqURL, "text/plain", http.NoBody)
}

func (c *Client) Update(m *model.Metric) error {
	if m.Name == "" {
		return errValidationName
	}
	reqURL, _ := url.JoinPath(c.address, updateEndpointPrefix)
	var bodyMessage handler.Metrics
	switch m.Type {
	case model.MetricTypeGauge:
		bodyMessage = handler.Metrics{
			ID:    m.Name,
			MType: m.Type.String(),
			Delta: nil,
			Value: &m.Gauge,
		}
	case model.MetricTypeCounter:
		bodyMessage = handler.Metrics{
			ID:    m.Name,
			MType: m.Type.String(),
			Delta: &m.Counter,
			Value: nil,
		}
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(bodyMessage)
	if err != nil {
		return err
	}

	return c.makeRequest(reqURL, "application/json", body)
}

func (c *Client) makeRequest(url string, contentType string, body io.Reader) error {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return errHTTP{Err: err}
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body error: %v, status code: %d", err, resp.StatusCode)
		}
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("close body error: %v, status code: %d", err, resp.StatusCode)
		}
		return errUnexpectedResponse{
			StatusCode: resp.StatusCode,
			Message:    body,
		}
	}
	return nil
}
