package client

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	gaugeEndpointPrefix   = "/update/gauge"
	counterEndpointPrefix = "/update/counter"
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

func (c *Client) Send(m internal.Metric) error {
	if m.Name == "" {
		return errValidationName
	}

	var reqURL string
	switch m.Type {
	case internal.GaugeType:
		reqURL, _ = url.JoinPath(
			c.address, gaugeEndpointPrefix, m.Name, strconv.FormatFloat(m.Gauge, 'f', -1, 64),
		)
	case internal.CounterType:
		reqURL, _ = url.JoinPath(
			c.address, counterEndpointPrefix, m.Name, fmt.Sprintf("%d", m.Counter),
		)
	}

	return c.makeRequest(reqURL)
}

func (c *Client) makeRequest(url string) error {
	resp, err := http.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return errHTTP{Err: err}
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("close body error: %v", err)
		}
		return errUnexpectedResponse{
			StatusCode: resp.StatusCode,
			Message:    body,
		}
	}
	return nil
}
