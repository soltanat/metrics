// Package client
// Клиент для отправки метрик
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/model"
)

const (
	gaugeEndpointPrefix   = "/update/gauge"
	counterEndpointPrefix = "/update/counter"
	updateEndpointPrefix  = "/update/"
	updatesEndpointPrefix = "/updates/"
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

// Client
// Клиент для отправки метрик
type Client struct {
	address string
	client  *http.Client
}

// New
// Создает новый клиент
// address - адрес API
func New(address string, transport http.RoundTripper) *Client {
	return &Client{
		address: address,
		client:  &http.Client{Transport: transport},
	}
}

// Send
// Отправляет метрику
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

// Update
// Отправляет метрику
func (c *Client) Update(m *model.Metric) error {
	if m.Name == "" {
		return errValidationName
	}

	reqURL, _ := url.JoinPath(c.address, updateEndpointPrefix)

	bodyMessage := handler.Metrics{
		ID:    m.Name,
		MType: m.Type.String(),
	}

	switch m.Type {
	case model.MetricTypeGauge:
		bodyMessage.Value = &m.Gauge
	case model.MetricTypeCounter:
		bodyMessage.Delta = &m.Counter
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(bodyMessage)
	if err != nil {
		return err
	}

	return c.makeRequest(reqURL, "application/json", body)
}

// Updates
// Обновляет слайс метрик
func (c *Client) Updates(metrics []model.Metric) error {
	reqURL, _ := url.JoinPath(c.address, updatesEndpointPrefix)

	bodyMessage := make([]handler.Metrics, 0, len(metrics))

	for i := 0; i < len(metrics); i++ {
		m := &metrics[i]
		bodyMetric := handler.Metrics{
			ID:    m.Name,
			MType: m.Type.String(),
		}
		switch m.Type {
		case model.MetricTypeGauge:
			bodyMetric.Value = &m.Gauge
		case model.MetricTypeCounter:
			bodyMetric.Delta = &m.Counter
		}
		bodyMessage = append(bodyMessage, bodyMetric)
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(bodyMessage)
	if err != nil {
		return err
	}

	return c.makeRequest(reqURL, "application/json", body)
}

func (c *Client) makeRequest(url string, contentType string, body io.Reader) error {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("create request error: %v", err)
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := c.client.Do(req)
	if err != nil {
		return errHTTP{Err: fmt.Errorf("request error: %v", err)}
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
