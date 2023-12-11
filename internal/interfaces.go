package internal

import "github.com/soltanat/metrics/internal/model"

type Poll interface {
	Get() ([]model.Metric, error)
	Poll() error
}

type Reporter interface {
	Report() error
}
