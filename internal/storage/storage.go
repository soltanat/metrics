package storage

import (
	"github.com/soltanat/metrics/internal/model"
)

// Storage
// Интерфейс хранилища метрик
type Storage interface {
	Store(metric *model.Metric) error
	StoreBatch(metrics []model.Metric) error
	GetGauge(name string) (*model.Metric, error)
	GetCounter(name string) (*model.Metric, error)
	GetList() ([]model.Metric, error)
}
