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

func (h *Handlers) GetGauge(c echo.Context) error {
	if err := validateMetricType(c); err != nil {
		return err
	}
	name := c.Param("metricName")
	metric, err := h.storage.GetGauge(name)
	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}
	_, _ = c.Response().Write([]byte(metric.AsString()))
	return nil
}

func (h *Handlers) GetCounter(c echo.Context) error {
	if err := validateMetricType(c); err != nil {
		return err
	}
	name := c.Param("metricName")
	metric, err := h.storage.GetCounter(name)
	if err != nil {
		if errors.Is(err, storage.ErrMetricNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}
	_, _ = c.Response().Write([]byte(metric.AsString()))
	return nil
}

func (h *Handlers) StoreCounter(c echo.Context) error {
	if err := validateMetricType(c); err != nil {
		return err
	}
	name := c.Param("metricName")
	value := c.Param("metricValue")

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}
	err = h.storage.StoreCounter(name, v)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (h *Handlers) StoreGauge(c echo.Context) error {
	if err := validateMetricType(c); err != nil {
		return err
	}
	name := c.Param("metricName")
	value := c.Param("metricValue")

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return echo.ErrBadRequest
	}
	err = h.storage.StoreGauge(name, v)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func validateMetricType(c echo.Context) error {
	mType := c.Param("metricType")
	switch mType {
	case internal.Gauge:
	case internal.Counter:
	default:
		return echo.ErrBadRequest
	}
	return nil
}
