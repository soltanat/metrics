package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/soltanat/metrics/internal/model"
)

const (
	gaugePathTmpl   = "/update/gauge/%s/%s"
	counterPathTmpl = "/update/counter/%s/%s"
)

var defaultMetricInst = model.NewGauge("name", 1.1)

func TestMetricsClient_Send(t *testing.T) {
	tests := []struct {
		name         string
		metric       *model.Metric
		expectedPath string
	}{
		{
			"gauge metric positive",
			model.NewGauge("name", 1.123),
			fmt.Sprintf(gaugePathTmpl, "name", "1.123"),
		},
		{
			"gauge metric round",
			model.NewGauge("name", 1),
			fmt.Sprintf(gaugePathTmpl, "name", "1"),
		},
		{
			"gauge metric negative",
			model.NewGauge("name", -1.123),
			fmt.Sprintf(gaugePathTmpl, "name", "-1.123"),
		},
		{
			"gauge metric zero",
			model.NewGauge("name", 0),
			fmt.Sprintf(gaugePathTmpl, "name", "0"),
		},
		{
			"counter metric zero",
			model.NewCounter("name", 0),
			fmt.Sprintf(counterPathTmpl, "name", "0"),
		},
		{
			"counter metric negative",
			model.NewCounter("name", -1),
			fmt.Sprintf(counterPathTmpl, "name", "-1"),
		},
		{
			"counter metric positive",
			model.NewCounter("name", 123),
			fmt.Sprintf(counterPathTmpl, "name", "123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.metric.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), tt.expectedPath)
				assert.Equal(t, req.Header.Get("content-type"), "text/plain")
				_, _ = rw.Write([]byte(`OK`))
			}))
			defer server.Close()

			api := Client{server.URL}
			err := api.Send(tt.metric)

			assert.NoError(t, err)
		})
	}
}

func TestMetricsClient_Send_IncorrectName(t *testing.T) {
	tests := []struct {
		name   string
		metric *model.Metric
	}{
		{
			"empty name",
			model.NewGauge("", 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := Client{"http://localhost:8080"}
			err := api.Send(tt.metric)
			assert.Error(t, err, errValidationName)
		})
	}
}

func TestMetricsClient_Send_ServerErrors(t *testing.T) {
	tests := []struct {
		name   string
		metric *model.Metric
	}{
		{
			"server error",
			defaultMetricInst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = rw.Write([]byte("error"))
			}))
			defer server.Close()

			api := Client{server.URL}
			err := api.Send(tt.metric)

			assert.Error(t, err, errUnexpectedResponse{
				StatusCode: 500,
				Message:    []byte("error"),
			})
		})
	}
}

func TestMetricsClient_Send_AddressError(t *testing.T) {
	tests := []struct {
		name   string
		metric *model.Metric
	}{
		{
			"address error",
			defaultMetricInst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := Client{"http://bad_address"}
			err := api.Send(tt.metric)

			var expectedErr errHTTP
			assert.ErrorAs(t, err, &expectedErr)
		})
	}
}

func TestUpdate(t *testing.T) {
	c := &Client{
		address: "http://example.com",
	}

	tests := []struct {
		name            string
		m               *model.Metric
		expectedErr     error
		expectedBody    string
		expectedReqBody string
	}{
		{
			name: "Valid Gauge Metric",
			m: &model.Metric{
				Name:  "metric1",
				Type:  model.MetricTypeGauge,
				Gauge: 1.5,
			},
			expectedErr:  nil,
			expectedBody: `{"id":"metric1","type":"gauge","value":1.5}` + "\n",
		},
		{
			name: "Valid Counter Metric",
			m: &model.Metric{
				Name:    "metric2",
				Type:    model.MetricTypeCounter,
				Counter: 10,
			},
			expectedErr:  nil,
			expectedBody: `{"id":"metric2","type":"counter","delta":10}` + "\n",
		},
		{
			name: "Invalid Metric (Empty Name)",
			m: &model.Metric{
				Name:  "",
				Type:  model.MetricTypeGauge,
				Gauge: 1.5,
			},
			expectedErr:  errValidationName,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				reqBody, _ := io.ReadAll(req.Body)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						t.Errorf("Error closing request body: %v", err)
					}
				}(req.Body)

				assert.Equal(t, tt.expectedBody, string(reqBody), "Expected request body: %s", tt.expectedReqBody)

				if tt.name == "Server Error" {
					rw.WriteHeader(http.StatusInternalServerError)
					_, _ = rw.Write([]byte("error"))
				}
			}))
			defer server.Close()

			c.address = server.URL

			err := c.Update(tt.m)
			if tt.expectedErr != nil {
				assert.Error(t, err, tt.expectedErr, "Expected error: %v", tt.expectedErr)
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}
