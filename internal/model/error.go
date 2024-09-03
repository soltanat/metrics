package model

import "fmt"

var (
	// ErrMetricNotFound
	// Ошибка при поиске метрики
	ErrMetricNotFound = fmt.Errorf("metric not found")
	ErrForbidden      = fmt.Errorf("forbidden")
)

type ErrBadRequest struct {
	Err error
}

func (e ErrBadRequest) Error() string {
	return e.Err.Error()
}

var (
	ErrMissingGaugeValue   = ErrBadRequest{Err: fmt.Errorf("missing gauge value")}
	ErrMissingCounterDelta = ErrBadRequest{Err: fmt.Errorf("missing counter delta")}
	ErrNotValidMetricType  = ErrBadRequest{Err: fmt.Errorf("not valid metric type")}
)
