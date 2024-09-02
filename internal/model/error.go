package model

import "fmt"

// ErrMetricNotFound
// Ошибка при поиске метрики
var ErrMetricNotFound = fmt.Errorf("metric not found")

var ErrForbidden = fmt.Errorf("forbidden")
