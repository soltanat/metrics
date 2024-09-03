package handler

import "github.com/soltanat/metrics/internal/model"

// Metrics схема передачи метрик
type Metrics struct {
	MID    string   `json:"id"`              // имя метрики
	MType  string   `json:"type"`            // параметр, принимающий значение gauge или counter
	MDelta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	MValue *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) Type() string {
	return m.MType
}

func (m Metrics) Value() *float64 {
	return m.MValue
}

func (m Metrics) Delta() *int64 {
	return m.MDelta
}

func (m Metrics) ID() string {
	return m.MID
}

func ConvertToInterfaces(input []Metrics) []model.InputMetric {
	metrics := make([]model.InputMetric, len(input))
	for i, m := range input {
		metrics[i] = m
	}
	return metrics
}
