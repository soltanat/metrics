package handler

import (
	"github.com/go-resty/resty/v2"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"
)

func TestHandlers_Get(t *testing.T) {
	type mockedFields struct {
		storage *storage.MockStorage
	}
	tests := []struct {
		name           string
		path           string
		mockedFields   mockedFields
		wantResponse   string
		wantStatusCode int
		on             func(fields *mockedFields)
	}{
		{
			name:           "existed counter",
			path:           "/value/counter/counter1",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}},
			wantResponse:   "1",
			wantStatusCode: http.StatusOK,
			on: func(fields *mockedFields) {
				fields.storage.On("GetCounter", "counter1").Return(
					model.NewCounter("counter1", 1), nil,
				)
			},
		},
		{
			name:           "existed gauge",
			path:           "/value/gauge/gauge1",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}},
			wantResponse:   "1.1",
			wantStatusCode: http.StatusOK,
			on: func(fields *mockedFields) {
				fields.storage.On("GetGauge", "gauge1").Return(
					model.NewGauge("gauge1", 1.1), nil,
				)
			},
		},
		{
			name:           "unknown type",
			path:           "/value/unknown/name",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "",
			on:             nil,
		},
		{
			name:           "not existed gauge",
			path:           "/value/gauge/name",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   "",
			on: func(fields *mockedFields) {
				fields.storage.On("GetGauge", "name").Return(
					nil, model.ErrMetricNotFound,
				)
			},
		},
		{
			name:           "not existed counter",
			path:           "/value/counter/name",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   "",
			on: func(fields *mockedFields) {
				fields.storage.On("GetCounter", "name").Return(
					nil, model.ErrMetricNotFound,
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.mockedFields.storage)
			srv := httptest.NewServer(SetupRoutes(h))
			defer srv.Close()

			req := resty.New().R()
			req.Method = resty.MethodGet

			u, err := url.JoinPath(srv.URL, tt.path)
			assert.NoError(t, err)

			req.URL = u

			if tt.on != nil {
				tt.on(&tt.mockedFields)
			}

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode())

			assert.True(t, tt.mockedFields.storage.AssertExpectations(t))
		})
	}
}

func TestHandlers_Store(t *testing.T) {
	type mockedFields struct {
		storage *storage.MockStorage
	}
	tests := []struct {
		name         string
		path         string
		mockedFields mockedFields
		wantRespCode int
		on           func(fields *mockedFields)
	}{
		{
			name:         "store counter",
			path:         "/update/counter/name/1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusOK,
			on: func(fields *mockedFields) {
				m := model.NewCounter("name", 1)
				fields.storage.On("GetCounter", "name").Return(m, nil)
				fields.storage.On("Store", m).Return(nil)
			},
		},
		{
			name:         "store new counter",
			path:         "/update/counter/name/1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusOK,
			on: func(fields *mockedFields) {
				m := model.NewCounter("name", 1)
				fields.storage.On("GetCounter", "name").Return(
					nil, model.ErrMetricNotFound,
				)
				fields.storage.On("Store", m).Return(nil)
			},
		},
		{
			name:         "store gauge",
			path:         "/update/gauge/name/1.1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusOK,
			on: func(fields *mockedFields) {
				fields.storage.On(
					"Store", model.NewGauge("name", 1.1),
				).Return(nil)
			},
		},
		{
			name:         "store float to counter",
			path:         "/update/counter/name/1.1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
		{
			name:         "store string to counter",
			path:         "/update/counter/name/str",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
		{
			name:         "store unknown type",
			path:         "/update/unknown/name/1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.mockedFields.storage)
			srv := httptest.NewServer(SetupRoutes(h))
			defer srv.Close()

			if tt.on != nil {
				tt.on(&tt.mockedFields)
			}

			req := resty.New().R()
			req.Method = resty.MethodPost
			u, err := url.JoinPath(srv.URL, tt.path)
			assert.NoError(t, err)
			req.URL = u

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantRespCode, resp.StatusCode())

			assert.True(t, tt.mockedFields.storage.AssertExpectations(t))
		})
	}
}

func TestHandlers_GetList(t *testing.T) {
	type mockedFields struct {
		storage *storage.MockStorage
	}
	tests := []struct {
		name           string
		mockedFields   mockedFields
		on             func(fields *mockedFields)
		wantedRespCode int
		wantedRespBody string
	}{
		{
			name:         "not empty storage",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			on: func(fields *mockedFields) {
				fields.storage.On("GetList").Return([]model.Metric{
					*model.NewCounter("metric1", 1),
					*model.NewGauge("metric2", 1.1),
				}, nil)
			},
			wantedRespCode: http.StatusOK,
			wantedRespBody: "type: counter, name: metric1, value: 1\ntype: gauge, name: metric2, value: 1.1\n",
		},
		{
			name:         "empty storage",
			mockedFields: mockedFields{storage: &storage.MockStorage{}},
			on: func(fields *mockedFields) {
				fields.storage.On("GetList").Return([]model.Metric{}, nil)
			},
			wantedRespCode: http.StatusOK,
			wantedRespBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.mockedFields.storage)
			srv := httptest.NewServer(SetupRoutes(h))
			defer srv.Close()

			tt.on(&tt.mockedFields)

			req := resty.New().R()
			req.Method = resty.MethodGet
			req.URL = srv.URL

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantedRespCode, resp.StatusCode())

			assert.True(t, tt.mockedFields.storage.AssertExpectations(t))

			assert.Equal(t, tt.wantedRespBody, string(resp.Body()))
		})
	}
}
