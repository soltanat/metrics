package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
	"strconv"
)

type Handlers struct {
	storage storage.Storage
}

func New(s storage.Storage) *Handlers {
	return &Handlers{storage: s}
}

func (h *Handlers) GetList(c echo.Context) error {
	metrics, err := h.storage.GetList()
	if err != nil {
		return echo.ErrInternalServerError
	}
	for _, m := range metrics {
		_, _ = c.Response().Write([]byte(m.AsString() + "\n"))
	}
	return nil
}

func (h *Handlers) Get(c echo.Context) error {
	metricTypeRaw := c.Param("metricType")
	name := c.Param("metricName")

	metricType, err := model.ParseMetricType(metricTypeRaw)
	if err != nil {
		log.Error(err)
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
	metricTypeRaw := c.Param("metricType")
	name := c.Param("metricName")
	valueRaw := c.Param("metricValue")

	var metric *model.Metric

	metricType, err := model.ParseMetricType(metricTypeRaw)
	if err != nil {
		log.Error(err)
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
		log.Error(err)
		return echo.ErrInternalServerError
	}

	return nil
}
