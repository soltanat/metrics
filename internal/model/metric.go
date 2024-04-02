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
