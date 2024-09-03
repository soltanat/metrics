package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/soltanat/metrics/internal/db"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
)

type Handlers struct {
	storage storage.Storage
	dbConn  db.Conn
	logger  zerolog.Logger
}

func New(s storage.Storage, dbConn db.Conn) *Handlers {
	return &Handlers{storage: s, dbConn: dbConn, logger: logger.Get()}
}

// GetList возвращает все метрики
func (h *Handlers) GetList(c echo.Context) error {
	metrics, err := h.storage.GetList()
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	for _, m := range metrics {
		_, _ = c.Response().Write([]byte(m.AsString() + "\n"))
	}

	return nil
}

// Get возвращает метрику
// metricType - тип метрики
// metricName - имя метрики
func (h *Handlers) Get(c echo.Context) error {
	l := logger.Get()

	metricTypeRaw := c.Param("metricType")
	name := c.Param("metricName")

	metricType, err := model.ParseMetricType(metricTypeRaw)
	if err != nil {
		l.Error().Err(err)
		return echo.ErrBadRequest
	}

	var metric *model.Metric

	switch metricType {
	case model.MetricTypeGauge:
		metric, err = h.storage.GetGauge(name)
	case model.MetricTypeCounter:
		metric, err = h.storage.GetCounter(name)
	}

	if err != nil {
		if errors.Is(err, model.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}
	_, _ = c.Response().Write([]byte(metric.ValueAsString()))

	return nil
}

// Store сохраняет метрику
// metricType - тип метрики
// metricName - имя метрики
// metricValue - значение метрики
func (h *Handlers) Store(c echo.Context) error {
	l := logger.Get()

	metricTypeRaw := c.Param("metricType")
	name := c.Param("metricName")
	valueRaw := c.Param("metricValue")

	var metric *model.Metric

	metricType, err := model.ParseMetricType(metricTypeRaw)
	if err != nil {
		l.Error().Err(err)
		return echo.ErrBadRequest
	}
	switch metricType {
	case model.MetricTypeGauge:
		value, parseErr := strconv.ParseFloat(valueRaw, 64)
		if parseErr != nil {
			return echo.ErrBadRequest
		}

		metric = model.NewGauge(name, value)

	case model.MetricTypeCounter:
		value, parseErr := strconv.ParseInt(valueRaw, 10, 64)
		if parseErr != nil {
			return echo.ErrBadRequest
		}
		metric = model.NewCounter(name, value)
	}

	err = h.storage.Store(metric)
	if err != nil {
		l.Error().Err(err)
		return echo.ErrInternalServerError
	}

	return nil
}

// StoreMetrics сохраняет метрику
// Тело запроса должно содержать JSON со схемой Metrics
func (h *Handlers) StoreMetrics(c echo.Context) error {
	var metrics Metrics
	if err := c.Bind(&metrics); err != nil {
		h.logger.Error().Msgf("Error binding metrics: %s", err)
		return echo.ErrBadRequest
	}

	metric, err := model.NewMetric(metrics)
	if err != nil {
		return err
	}

	if err := h.storage.Store(metric); err != nil {
		h.logger.Error().Msgf("Error storing metric: %s", err)
		return echo.ErrInternalServerError
	}

	h.logger.Info().Msg("Metrics stored successfully")

	return c.NoContent(http.StatusOK)
}

// StoreMetricsBatch сохраняет метрики
// Тело запроса должно содержать JSON список элементов Metrics
func (h *Handlers) StoreMetricsBatch(c echo.Context) error {
	var metrics []Metrics
	if err := c.Bind(&metrics); err != nil {
		h.logger.Error().Msgf("Error binding metrics: %s", err)
		return echo.ErrBadRequest
	}

	mm, err := model.NewMetrics(ConvertToInterfaces(metrics))
	if err != nil {
		if errors.As(err, &model.ErrBadRequest{}) {
			h.logger.Error().Err(err).Msg("Error parsing metrics")
			return echo.ErrBadRequest
		}
	}

	if err := h.storage.StoreBatch(mm); err != nil {
		h.logger.Error().Msgf("Error storing metric: %s", err)
		return echo.ErrInternalServerError
	}

	h.logger.Info().Msg("Metrics stored successfully")

	return c.NoContent(http.StatusOK)
}

// Value возвращает значение метрики
func (h *Handlers) Value(c echo.Context) error {
	var m Metrics
	if err := c.Bind(&m); err != nil {
		return echo.ErrBadRequest
	}

	metricType, err := model.ParseMetricType(m.MType)
	if err != nil {
		return echo.ErrBadRequest
	}

	var metric *model.Metric
	metricMap := map[model.MetricType]func(string) (*model.Metric, error){
		model.MetricTypeGauge:   h.storage.GetGauge,
		model.MetricTypeCounter: h.storage.GetCounter,
	}

	metricFunc, ok := metricMap[metricType]
	if !ok {
		return echo.ErrBadRequest
	}

	metric, err = metricFunc(m.MID)
	if err != nil {
		if errors.Is(err, model.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	switch metricType {
	case model.MetricTypeGauge:
		m.MValue = &metric.Gauge
	case model.MetricTypeCounter:
		m.MDelta = &metric.Counter
	}

	return c.JSON(http.StatusOK, m)
}

// Ping проверяет соединение с базой данных
// Если соединение установлено возвращает 200
// Если соединение не установлено возвращает 503
func (h *Handlers) Ping(c echo.Context) error {
	if h.dbConn == nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	if err := h.dbConn.Ping(c.Request().Context()); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
