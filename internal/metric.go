package internal

import (
	"fmt"
	"strconv"
)

type MetricType int

const (
	GaugeType MetricType = iota
	CounterType
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric struct {
	Type    MetricType
	Name    string
	Gauge   float64
	Counter int64
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
	case GaugeType:
		v := strconv.FormatFloat(m.Gauge, 'f', -1, 64)
		return fmt.Sprintf("type: %s, name: %s, value: %s", Gauge, m.Name, v)
	case CounterType:
		return fmt.Sprintf("type: %s, name: %s, value: %d", Counter, m.Name, m.Counter)
	}
	return ""
}

func (m *Metric) ValueAsString() string {
	switch m.Type {
	case GaugeType:
		v := strconv.FormatFloat(m.Gauge, 'f', -1, 64)
		return v
	case CounterType:
		return fmt.Sprintf("%d", m.Counter)
	}
	return ""
}
