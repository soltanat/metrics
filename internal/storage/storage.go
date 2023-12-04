package storage

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Storage interface {
	StoreGauge(name string, value float64) error
	StoreCounter(name string, value int64) error
}
