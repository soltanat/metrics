package handler

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/soltanat/metrics/internal/db"
	"github.com/soltanat/metrics/internal/db/mock"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
)

func TestHandlers_Get(t *testing.T) {
	type mockedFields struct {
		storage *storage.MockStorage
		db      db.Conn
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
			mockedFields:   mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			mockedFields:   mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			mockedFields:   mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   "",
			on:             nil,
		},
		{
			name:           "not existed gauge",
			path:           "/value/gauge/name",
			mockedFields:   mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			mockedFields:   mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			h := New(tt.mockedFields.storage, tt.mockedFields.db)
			r, err := SetupRoutes(h, "", []byte(""), "")
			require.NoError(t, err)
			srv := httptest.NewServer(r)
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
		db      db.Conn
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
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantRespCode: http.StatusOK,
			on: func(fields *mockedFields) {
				m := model.NewCounter("name", 1)
				fields.storage.On("Store", m).Return(nil)
			},
		},
		{
			name:         "store new counter",
			path:         "/update/counter/name/1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantRespCode: http.StatusOK,
			on: func(fields *mockedFields) {
				m := model.NewCounter("name", 1)
				fields.storage.On("Store", m).Return(nil)
			},
		},
		{
			name:         "store gauge",
			path:         "/update/gauge/name/1.1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
		{
			name:         "store string to counter",
			path:         "/update/counter/name/str",
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
		{
			name:         "store unknown type",
			path:         "/update/unknown/name/1",
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			wantRespCode: http.StatusBadRequest,
			on:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.mockedFields.storage, tt.mockedFields.db)
			r, err := SetupRoutes(h, "", []byte(""), "")
			require.NoError(t, err)
			srv := httptest.NewServer(r)
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
		db      db.Conn
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
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
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
			mockedFields: mockedFields{storage: &storage.MockStorage{}, db: &mock.MockConn{}},
			on: func(fields *mockedFields) {
				fields.storage.On("GetList").Return([]model.Metric{}, nil)
			},
			wantedRespCode: http.StatusOK,
			wantedRespBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.mockedFields.storage, tt.mockedFields.db)
			r, err := SetupRoutes(h, "", []byte(""), "")
			require.NoError(t, err)
			srv := httptest.NewServer(r)
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

type StorageCall struct {
	Metric *model.Metric
}

func TestHandlers_StoreMetrics(t *testing.T) {
	mockStorage := &storage.MockStorage{}

	h := &Handlers{
		logger:  logger.Get(),
		storage: mockStorage,
	}

	r, err := SetupRoutes(h, "", []byte(""), "")
	require.NoError(t, err)
	server := httptest.NewServer(r)

	defer server.Close()

	client := resty.New()

	testCases := []struct {
		name         string
		metrics      Metrics
		expectedCall *StorageCall
		statusCode   int
		on           func(metric *model.Metric, storage *storage.MockStorage)
	}{
		{
			name: "StoreGaugeMetric",
			metrics: Metrics{
				MType:  "gauge",
				MID:    "test-id",
				MValue: float64Ptr(10.1),
			},
			expectedCall: &StorageCall{
				Metric: &model.Metric{
					Type:  model.MetricTypeGauge,
					Name:  "test-id",
					Gauge: 10.1,
				},
			},
			statusCode: http.StatusOK,
			on: func(metric *model.Metric, storage *storage.MockStorage) {
				mockStorage.On("Store", metric).Return(nil).Once()
			},
		},
		{
			name: "StoreCounterMetricExistMetric",
			metrics: Metrics{
				MType:  "counter",
				MID:    "test-id",
				MDelta: intPtr(5),
			},
			expectedCall: &StorageCall{
				Metric: &model.Metric{
					Type:    model.MetricTypeCounter,
					Name:    "test-id",
					Counter: 5,
				},
			},
			statusCode: http.StatusOK,
			on: func(metric *model.Metric, storage *storage.MockStorage) {
				mockStorage.On("GetCounter", "test-id").Return(metric, nil).Once()
				mockStorage.On("Store", metric).Return(nil).Once()
			},
		},
		{
			name: "StoreCounterMetricNotExistMetric",
			metrics: Metrics{
				MType:  "counter",
				MID:    "test-id",
				MDelta: intPtr(5),
			},
			expectedCall: &StorageCall{
				Metric: &model.Metric{
					Type:    model.MetricTypeCounter,
					Name:    "test-id",
					Counter: 5,
				},
			},
			statusCode: http.StatusOK,
			on: func(metric *model.Metric, storage *storage.MockStorage) {
				mockStorage.On("GetCounter", "test-id").Return(
					nil,
					model.ErrMetricNotFound).Once()
				mockStorage.On("Store", metric).Return(nil).Once()
			},
		},
		{
			name: "MissingValueForGaugeMetric",
			metrics: Metrics{
				MType: "gauge",
				MID:   "test-id",
			},
			expectedCall: nil,
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "MissingDeltaForCounterMetric",
			metrics: Metrics{
				MType: "counter",
				MID:   "test-id",
			},
			expectedCall: nil,
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "UnknownMetricType",
			metrics: Metrics{
				MType: "unknown",
				MID:   "test-id",
			},
			expectedCall: nil,
			statusCode:   http.StatusBadRequest,
		},
		{
			name: "InvalidRequestBody",
			metrics: Metrics{
				MType: "gauge",
				MID:   "test-id",
			},
			expectedCall: nil,
			statusCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.on != nil {
				tc.on(tc.expectedCall.Metric, mockStorage)
			}

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tc.metrics).
				Post(server.URL + "/update/")
			assert.NoError(t, err)
			assert.Equal(t, tc.statusCode, resp.StatusCode())

			if tc.expectedCall != nil {
				mockStorage.AssertCalled(t, "Store", tc.expectedCall.Metric)
			}
		})
	}
}

func TestHandlers_Value(t *testing.T) {
	mockStorage := &storage.MockStorage{}
	h := &Handlers{
		storage: mockStorage,
	}

	r, err := SetupRoutes(h, "", []byte(""), "")
	require.NoError(t, err)
	server := httptest.NewServer(r)

	defer server.Close()

	client := resty.New()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		on             func()
	}{
		{
			name:           "Invalid request body",
			requestBody:    `invalid request`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"message\":\"Bad Request\"}",
			on:             nil,
		},
		{
			name:           "Invalid metric type",
			requestBody:    `{"type": "invalid_type"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"message\":\"Bad Request\"}",
		},
		{
			name:           "Metric not found",
			requestBody:    `{"type": "gauge", "id": "not_found_id"}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "{\"message\":\"Not Found\"}",
			on: func() {
				mockStorage.On("GetGauge", "not_found_id").Return(
					nil, model.ErrMetricNotFound,
				).Once()
			},
		},
		{
			name:           "Successful request - Gauge",
			requestBody:    `{"type": "gauge", "id": "valid_id"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"type": "gauge", "id": "valid_id", "value": 10.1}`,
			on: func() {
				mockStorage.On("GetGauge", "valid_id").Return(
					model.NewGauge("valid_id", 10.1), nil,
				).Once()
			},
		},
		{
			name:           "Successful request - Counter",
			requestBody:    `{"type": "counter", "id": "valid_id"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"type": "counter", "id": "valid_id", "delta": 10}`,
			on: func() {
				mockStorage.On("GetCounter", "valid_id").Return(
					model.NewCounter("valid_id", 10), nil,
				).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.on != nil {
				test.on()
			}

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(test.requestBody).
				Post(server.URL + "/value/")

			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatus, resp.StatusCode())
			assert.JSONEq(t, test.expectedBody, resp.String())

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestHandlers_Ping(t *testing.T) {
	mockStorage := &storage.MockStorage{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockConn(ctrl)

	h := &Handlers{
		storage: mockStorage,
		dbConn:  mockDB,
	}

	r, err := SetupRoutes(h, "", []byte(""), "")
	require.NoError(t, err)
	server := httptest.NewServer(r)

	defer server.Close()

	client := resty.New()

	tests := []struct {
		name           string
		expectedStatus int
		on             func()
	}{
		{
			name:           "Ping failed",
			expectedStatus: http.StatusInternalServerError,
			on: func() {
				mockDB.EXPECT().Ping(gomock.Any()).Return(fmt.Errorf("ping failed")).Times(1)
			},
		},
		{
			name:           "Ping succeeded",
			expectedStatus: http.StatusOK,
			on: func() {
				mockDB.EXPECT().Ping(gomock.Any()).Return(nil).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.on != nil {
				test.on()
			}

			resp, err := client.R().Get(server.URL + "/ping")

			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatus, resp.StatusCode())
			assert.Equal(t, "", resp.String())

			mockStorage.AssertExpectations(t)
		})
	}
}

// Helper function to create a float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}

// Helper function to create an int pointer
func intPtr(i int64) *int64 {
	return &i
}
