package model

import (
	"fmt"
	"strconv"
)

// Metric
// Структура метрики
// type: тип метрики (gauge, counter)
// name: имя метрики
// gauge: значение gauge
// counter: значение counter
type Metric struct {
	Type    MetricType
	Name    string
	Gauge   float64
	Counter int64
}

type InputMetric interface {
	Type() string
	Value() *float64
	Delta() *int64
	ID() string
}

func NewMetrics(input []InputMetric) ([]Metric, error) {
	metrics := make([]Metric, 0, len(input))
	for _, metric := range input {
		m, err := NewMetric(metric)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, *m)
	}
	return metrics, nil
}

func NewMetric(input InputMetric) (*Metric, error) {
	var metric *Metric

	metricType, err := ParseMetricType(input.Type())
	if err != nil {
		return nil, ErrNotValidMetricType
	}

	switch metricType {
	case MetricTypeGauge:
		if input.Value() == nil {
			return nil, ErrMissingGaugeValue
		}
		metric = NewGauge(input.ID(), *input.Value())
	case MetricTypeCounter:
		if input.Delta() == nil {
			return nil, ErrMissingCounterDelta
		}
		metric = NewCounter(input.ID(), *input.Delta())
	default:
		return nil, ErrNotValidMetricType
	}

	return metric, nil
}

func NewGauge(name string, value float64) *Metric {
	return &Metric{
		Type:  MetricTypeGauge,
		Name:  name,
		Gauge: value,
	}
}

func NewCounter(name string, value int64) *Metric {
	return &Metric{
		Type:    MetricTypeCounter,
		Name:    name,
		Counter: value,
	}
}

// IncCounter
// Увеличивает значение counter
func (m *Metric) IncCounter() {
	m.Counter += 1
}

// SetGauge
// Устанавливает значение gauge
func (m *Metric) SetGauge(v float64) {
	m.Gauge = v
}

// AsString
// Возвращает строковое представление метрики
func (m *Metric) AsString() string {
	switch m.Type {
	case MetricTypeGauge:
		v := strconv.FormatFloat(m.Gauge, 'f', -1, 64)
		return fmt.Sprintf("type: %s, name: %s, value: %s", MetricTypeGauge.String(), m.Name, v)
	case MetricTypeCounter:
		return fmt.Sprintf("type: %s, name: %s, value: %d", MetricTypeCounter.String(), m.Name, m.Counter)
	}
	return ""
}

// ValueAsString
// Возвращает строковое представление значения
func (m *Metric) ValueAsString() string {
	switch m.Type {
	case MetricTypeGauge:
		v := strconv.FormatFloat(m.Gauge, 'f', -1, 64)
		return v
	case MetricTypeCounter:
		return fmt.Sprintf("%d", m.Counter)
	}
	return ""
}
