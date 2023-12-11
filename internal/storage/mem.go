package storage

import (
	"github.com/soltanat/metrics/internal/model"
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      *sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
		mu:      &sync.RWMutex{},
	}
}

func (s *MemStorage) Store(metric *model.Metric) error {
	switch metric.Type {
	case model.MetricTypeCounter:
		s.counter[metric.Name] = metric.Counter
	case model.MetricTypeGauge:
		s.gauge[metric.Name] = metric.Gauge
	}
	return nil
}

func (s *MemStorage) GetGauge(name string) (*model.Metric, error) {
	v, ok := s.gauge[name]
	if !ok {
		return nil, model.ErrMetricNotFound
	}
	return model.NewGauge(name, v), nil
}

func (s *MemStorage) GetCounter(name string) (*model.Metric, error) {
	v, ok := s.counter[name]
	if !ok {
		return nil, model.ErrMetricNotFound
	}
	return model.NewCounter(name, v), nil
}

func (s *MemStorage) GetList() ([]model.Metric, error) {
	var metrics []model.Metric
	for k, v := range s.counter {
		metrics = append(metrics, *model.NewCounter(k, v))
	}
	for k, v := range s.gauge {
		metrics = append(metrics, *model.NewGauge(k, v))
	}
	return metrics, nil
}
