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

var validationNameError = fmt.Errorf("min name len 1")

type HTTPError struct {
	Err error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

type UnexpectedResponseError struct {
	StatusCode int
	Message    []byte
}

func (e UnexpectedResponseError) Error() string {
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
		return validationNameError
	}

	var reqUrl string
	switch m.Type {
	case internal.GaugeType:
		reqUrl, _ = url.JoinPath(
			c.address, gaugeEndpointPrefix, m.Name, strconv.FormatFloat(m.Gauge, 'f', -1, 64),
		)
	case internal.CounterType:
		reqUrl, _ = url.JoinPath(
			c.address, counterEndpointPrefix, m.Name, fmt.Sprintf("%d", m.Counter),
		)
	}

	return c.makeRequest(reqUrl)
}

func (c *Client) makeRequest(url string) error {
	resp, err := http.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return HTTPError{Err: err}
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
		return UnexpectedResponseError{
			StatusCode: resp.StatusCode,
			Message:    body,
		}
	}
	return nil
}
