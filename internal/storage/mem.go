package storage

import "github.com/soltanat/metrics/internal"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *MemStorage) Store(metric *internal.Metric) error {
	switch metric.Type {
	case internal.CounterType:
		s.counter[metric.Name] = metric.Counter
	case internal.GaugeType:
		s.gauge[metric.Name] = metric.Gauge
	}
	return nil
}

func (s *MemStorage) GetGauge(name string) (*internal.Metric, error) {
	m, ok := s.gauge[name]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return &internal.Metric{
		Type:    internal.GaugeType,
		Name:    name,
		Gauge:   m,
		Counter: 0,
	}, nil
}

func (s *MemStorage) GetCounter(name string) (*internal.Metric, error) {
	m, ok := s.counter[name]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return &internal.Metric{
		Type:    internal.CounterType,
		Name:    name,
		Gauge:   0,
		Counter: m,
	}, nil
}

func (s *MemStorage) GetList() ([]internal.Metric, error) {
	var metrics []internal.Metric
	for k, v := range s.counter {
		metrics = append(metrics, internal.Metric{
			Type:    internal.CounterType,
			Name:    k,
			Gauge:   0,
			Counter: v,
		})
	}
	for k, v := range s.gauge {
		metrics = append(metrics, internal.Metric{
			Type:    internal.GaugeType,
			Name:    k,
			Gauge:   v,
			Counter: 0,
		})
	}
	return metrics, nil
}
