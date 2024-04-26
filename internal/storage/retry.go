package storage

import (
	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/model"
)

// BackoffPostgresStorage
// Декоратор Storage с попытками повтороной обработки ошибок
type BackoffPostgresStorage struct {
	storage Storage
}

func NewBackoffPostgresStorage(s Storage) *BackoffPostgresStorage {
	return &BackoffPostgresStorage{storage: s}
}

func (s *BackoffPostgresStorage) Store(metric *model.Metric) error {
	return internal.Backoff(func() error {
		return s.storage.Store(metric)
	})
}

func (s *BackoffPostgresStorage) StoreBatch(metrics []model.Metric) error {
	return internal.Backoff(func() error {
		return s.storage.StoreBatch(metrics)
	})
}

func (s *BackoffPostgresStorage) GetGauge(name string) (metric *model.Metric, err error) {
	err = internal.Backoff(func() error {
		metric, err = s.storage.GetGauge(name)
		return err
	}, model.ErrMetricNotFound)
	return
}

func (s *BackoffPostgresStorage) GetCounter(name string) (metric *model.Metric, err error) {
	err = internal.Backoff(func() error {
		metric, err = s.storage.GetCounter(name)
		return err
	}, model.ErrMetricNotFound)
	return
}

func (s *BackoffPostgresStorage) GetList() (metrics []model.Metric, err error) {
	err = internal.Backoff(func() error {
		metrics, err = s.storage.GetList()
		return err
	})
	return
}
