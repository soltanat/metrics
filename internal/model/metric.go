package model

import (
	"fmt"
	"strconv"
)

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

func (m *Metric) IncCounter() {
	m.Counter += 1
}

func (m *Metric) AddCounter(v int64) {
	m.Counter += v
}

func (m *Metric) SetGauge(v float64) {
	m.Gauge = v
}

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
