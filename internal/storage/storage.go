package storage

import (
	"github.com/soltanat/metrics/internal/model"
)

type Storage interface {
	Store(metric *model.Metric) error
	GetGauge(name string) (*model.Metric, error)
	GetCounter(name string) (*model.Metric, error)
	GetList() ([]model.Metric, error)
}
