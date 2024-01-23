package handler

import (
	"errors"
	"github.com/soltanat/metrics/internal/db"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
)

type Handlers struct {
	storage storage.Storage
	db      *db.DB
	logger  zerolog.Logger
}

func New(s storage.Storage, db *db.DB) *Handlers {
	return &Handlers{storage: s, db: db, logger: logger.Get()}
}

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
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		metric = model.NewGauge(name, value)

	case model.MetricTypeCounter:
		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}

		metric, err = h.storage.GetCounter(name)

		if err != nil {
			if !errors.Is(err, model.ErrMetricNotFound) {
				return echo.ErrBadRequest
			}
			metric = model.NewCounter(name, 0)
		}

		metric.AddCounter(value)
	}

	err = h.storage.Store(metric)
	if err != nil {
		l.Error().Err(err)
		return echo.ErrInternalServerError
	}

	return nil
}

func (h *Handlers) StoreMetrics(c echo.Context) error {
	var metrics Metrics
	if err := c.Bind(&metrics); err != nil {
		h.logger.Error().Msgf("Error binding metrics: %s", err)
		return echo.ErrBadRequest
	}

	var metric *model.Metric

	metricType, err := model.ParseMetricType(metrics.MType)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error parsing metric type")
		return echo.ErrBadRequest
	}

	switch metricType {
	case model.MetricTypeGauge:
		if metrics.Value == nil {
			h.logger.Error().Msg("Missing value for gauge metric")
			return echo.ErrBadRequest
		}
		metric = model.NewGauge(metrics.ID, *metrics.Value)
	case model.MetricTypeCounter:
		if metrics.Delta == nil {
			h.logger.Error().Msg("Missing delta for counter metric")
			return echo.ErrBadRequest
		}
		m, err := h.storage.GetCounter(metrics.ID)
		if err != nil {
			if !errors.Is(err, model.ErrMetricNotFound) {
				h.logger.Error().Msgf("Error getting counter metric: %s", err)
				return echo.ErrBadRequest
			}
			m = model.NewCounter(metrics.ID, 0)
		}
		m.AddCounter(*metrics.Delta)
		metric = m
	default:
		h.logger.Error().Msg("Unknown metric type")
		return echo.ErrBadRequest
	}

	if err := h.storage.Store(metric); err != nil {
		h.logger.Error().Msgf("Error storing metric: %s", err)
		return echo.ErrInternalServerError
	}

	h.logger.Info().Msg("Metrics stored successfully")

	return c.NoContent(http.StatusOK)
}

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

	metric, err = metricFunc(m.ID)
	if err != nil {
		if errors.Is(err, model.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	switch metricType {
	case model.MetricTypeGauge:
		m.Value = &metric.Gauge
	case model.MetricTypeCounter:
		m.Delta = &metric.Counter
	}

	return c.JSON(http.StatusOK, m)
}

func (h *Handlers) Ping(c echo.Context) error {
	if err := h.db.Ping(c.Request().Context()); err != nil {
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}
