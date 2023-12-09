package internal

type MetricType int

const (
	GaugeType MetricType = iota
	CounterType
)

type Metric struct {
	Type    MetricType
	Name    string
	Gauge   float64
	Counter int64
}

func (m *Metric) IncCount() {
	m.Counter += 1
}

func (m *Metric) SetGauge(v float64) {
	m.Gauge = v
}
