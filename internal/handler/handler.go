package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/soltanat/metrics/internal"
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
	mType := c.Param("metricType")
	name := c.Param("metricName")

	var metric *internal.Metric
	var err error

	switch mType {
	case internal.Gauge:
		metric, err = h.storage.GetGauge(name)
	case internal.Counter:
		metric, err = h.storage.GetCounter(name)
	default:
		return echo.ErrBadRequest
	}

	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}
	_, _ = c.Response().Write([]byte(metric.ValueAsString()))

	return nil

}

func (h *Handlers) Store(c echo.Context) error {
	mType := c.Param("metricType")
	name := c.Param("metricName")
	value := c.Param("metricValue")

	switch mType {
	case internal.Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return echo.ErrBadRequest
		}
		err = h.storage.StoreGauge(name, v)
		if err != nil {
			return echo.ErrBadRequest
		}
	case internal.Counter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return echo.ErrBadRequest
		}
		err = h.storage.StoreCounter(name, v)
		if err != nil {
			return echo.ErrBadRequest
		}
	default:
		return echo.ErrBadRequest
	}

	return nil
}
