package storage

import (
	"sync"

	"github.com/soltanat/metrics/internal/model"
)

// MemStorage
// Реализует хранилище метрик в памяти
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

// Store
// Сохраняет метрику
// Для counter добавляет значение, для gauge заменяет значение
func (s *MemStorage) Store(metric *model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store(metric)
}

// StoreBatch
// Сохраняет слайс метрик
// Для counter добавляет значения, для gauge заменяет значения
func (s *MemStorage) StoreBatch(metrics []model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, metric := range metrics {
		if err := s.store(&metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemStorage) store(metric *model.Metric) error {
	switch metric.Type {
	case model.MetricTypeCounter:
		s.counter[metric.Name] += metric.Counter
	case model.MetricTypeGauge:
		s.gauge[metric.Name] = metric.Gauge
	}
	return nil
}

// GetGauge
// Возвращает метрику gauge по имени
func (s *MemStorage) GetGauge(name string) (*model.Metric, error) {
	s.mu.RLock()
	v, ok := s.gauge[name]
	s.mu.RUnlock()
	if !ok {
		return nil, model.ErrMetricNotFound
	}
	return model.NewGauge(name, v), nil
}

// GetCounter
// Возвращает метрику counter по имени
func (s *MemStorage) GetCounter(name string) (*model.Metric, error) {
	s.mu.RLock()
	v, ok := s.counter[name]
	s.mu.RUnlock()
	if !ok {
		return nil, model.ErrMetricNotFound
	}
	return model.NewCounter(name, v), nil
}

// GetList
// Возвращает все метрики в виде слайса
func (s *MemStorage) GetList() ([]model.Metric, error) {
	s.mu.RLock()
	metrics := make([]model.Metric, 0, len(s.counter)+len(s.gauge))
	for k, v := range s.counter {
		metrics = append(metrics, *model.NewCounter(k, v))
	}
	for k, v := range s.gauge {
		metrics = append(metrics, *model.NewGauge(k, v))
	}
	s.mu.RUnlock()
	return metrics, nil
}
