package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
)

var l = logger.Get()

func ExampleHandlers_GetList() {

	s := storage.NewMemStorage()
	_ = s.Store(&model.Metric{
		Name:  "test",
		Type:  model.MetricTypeGauge,
		Gauge: 1,
	})
	_ = s.Store(&model.Metric{
		Name:    "test",
		Type:    model.MetricTypeCounter,
		Counter: 1,
	})

	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/`, server.URL), http.NoBody)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// type: counter, name: test, value: 1
	// type: gauge, name: test, value: 1
}

func ExampleHandlers_Store() {
	s := storage.NewMemStorage()
	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(`%s/update/gauge/test/1/`, server.URL), http.NoBody)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
}

func ExampleHandlers_StoreMetricsBatch() {
	s := storage.NewMemStorage()
	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	metrics := []Metrics{
		{
			MID:    "test",
			MType:  "gauge",
			MValue: float64Ptr(10.1),
		},
		{
			MID:    "test",
			MType:  "counter",
			MDelta: intPtr(10),
		},
	}

	reqBody, err := json.Marshal(metrics)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to marshal metrics")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(`%s/updates/`, server.URL), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
}

func ExampleHandlers_StoreMetrics() {
	s := storage.NewMemStorage()
	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	metrics := Metrics{
		MID:    "test",
		MType:  "gauge",
		MValue: float64Ptr(10.1),
	}

	reqBody, err := json.Marshal(metrics)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to marshal metrics")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(`%s/update/`, server.URL), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
}

func ExampleHandlers_Value() {
	s := storage.NewMemStorage()
	_ = s.Store(&model.Metric{
		Name:  "test",
		Type:  model.MetricTypeGauge,
		Gauge: 1,
	})
	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	metrics := Metrics{
		MID:   "test",
		MType: "gauge",
	}

	reqBody, err := json.Marshal(metrics)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to marshal metrics")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(`%s/value/`, server.URL), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// {"id":"test","type":"gauge","value":1}
}

func ExampleHandlers_Get() {
	s := storage.NewMemStorage()
	_ = s.Store(&model.Metric{
		Name:  "test",
		Type:  model.MetricTypeGauge,
		Gauge: 1,
	})
	h := New(s, nil)

	r, err := SetupRoutes(h, "", []byte(""))
	if err != nil {
		l.Fatal().Err(err).Msg("failed to setup routes")
	}
	server := httptest.NewServer(r)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/value/gauge/test`, server.URL), nil)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to do request")
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// 1
}
