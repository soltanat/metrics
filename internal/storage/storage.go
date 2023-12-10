package storage

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
)

var ErrMetricNotFound = fmt.Errorf("metric not found")

type Storage interface {
	StoreGauge(name string, value float64) error
	StoreCounter(name string, value int64) error
	GetGauge(name string) (*internal.Metric, error)
	GetCounter(name string) (*internal.Metric, error)
	GetList() ([]internal.Metric, error)
}
