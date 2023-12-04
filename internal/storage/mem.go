package storage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s MemStorage) StoreGauge(name string, value float64) error {
	s.gauge[name] = value
	return nil
}

func (s MemStorage) StoreCounter(name string, value int64) error {
	s.counter[name] += value
	return nil
}
