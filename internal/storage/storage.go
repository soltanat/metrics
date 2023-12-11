package storage

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
)

var ErrMetricNotFound = fmt.Errorf("metric not found")

type Storage interface {
	Store(metric *internal.Metric) error
	GetGauge(name string) (*internal.Metric, error)
	GetCounter(name string) (*internal.Metric, error)
	GetList() ([]internal.Metric, error)
}
