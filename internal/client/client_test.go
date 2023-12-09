package client

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	gaugePathTmpl   = "/update/gauge/%s/%s"
	counterPathTmpl = "/update/counter/%s/%s"
)

var defaultMetricInst = internal.Metric{
	Type:    internal.GaugeType,
	Name:    "name",
	Counter: 1,
	Gauge:   0,
}

func TestMetricsClient_Send(t *testing.T) {
	tests := []struct {
		name         string
		metric       internal.Metric
		expectedPath string
	}{
		{
			"gauge metric positive",
			internal.Metric{
				Type:    internal.GaugeType,
				Name:    "name",
				Counter: 0,
				Gauge:   1.123,
			},
			fmt.Sprintf(gaugePathTmpl, "name", "1.123"),
		},
		{
			"gauge metric round",
			internal.Metric{
				Type:    internal.GaugeType,
				Name:    "name",
				Counter: 0,
				Gauge:   1,
			},
			fmt.Sprintf(gaugePathTmpl, "name", "1"),
		},
		{
			"gauge metric negative",
			internal.Metric{
				Type:    internal.GaugeType,
				Name:    "name",
				Counter: 0,
				Gauge:   -1.123,
			},
			fmt.Sprintf(gaugePathTmpl, "name", "-1.123"),
		},
		{
			"gauge metric zero",
			internal.Metric{
				Type:    internal.GaugeType,
				Name:    "name",
				Counter: 0,
				Gauge:   0,
			},
			fmt.Sprintf(gaugePathTmpl, "name", "0"),
		},
		{
			"counter metric zero",
			internal.Metric{
				Type:    internal.CounterType,
				Name:    "name",
				Counter: 0,
				Gauge:   0,
			},
			fmt.Sprintf(counterPathTmpl, "name", "0"),
		},
		{
			"counter metric negative",
			internal.Metric{
				Type:    internal.CounterType,
				Name:    "name",
				Counter: -1,
				Gauge:   0,
			},
			fmt.Sprintf(counterPathTmpl, "name", "-1"),
		},
		{
			"counter metric positive",
			internal.Metric{
				Type:    internal.CounterType,
				Name:    "name",
				Counter: 123,
				Gauge:   0,
			},
			fmt.Sprintf(counterPathTmpl, "name", "123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.metric.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), tt.expectedPath)
				assert.Equal(t, req.Header.Get("content-type"), "text/plain")
				rw.Write([]byte(`OK`))
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
		metric internal.Metric
	}{
		{
			"empty name",
			internal.Metric{
				Type:    internal.GaugeType,
				Name:    "",
				Counter: 0,
				Gauge:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := Client{"http://localhost:8080"}
			err := api.Send(tt.metric)
			assert.Error(t, err, validationNameError)
		})
	}
}

func TestMetricsClient_Send_ServerErrors(t *testing.T) {
	tests := []struct {
		name   string
		metric internal.Metric
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
				rw.Write([]byte("error"))
			}))
			defer server.Close()

			api := Client{server.URL}
			err := api.Send(tt.metric)

			assert.Error(t, err, UnexpectedResponseError{
				StatusCode: 500,
				Message:    []byte("error"),
			})
		})
	}
}

func TestMetricsClient_Send_AddressError(t *testing.T) {
	tests := []struct {
		name   string
		metric internal.Metric
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

			var expectedErr HTTPError
			assert.ErrorAs(t, err, &expectedErr)
		})
	}
}
