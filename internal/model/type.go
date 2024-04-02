package model

//go:generate go-enum --marshal

// MetricType ENUM(gauge, counter)
// тип метрики
// gauge: тип метрики gauge
// counter: тип метрики counter
type MetricType int
