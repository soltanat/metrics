package internal

import (
	"context"
	"github.com/soltanat/metrics/internal/model"
	"time"
)

type Poll interface {
	RunPoller(ctx context.Context, interval time.Duration) error
	GetChannel() chan *model.Metric
}

type Reporter interface {
	RunReporter(ctx context.Context, interval time.Duration, ch chan *model.Metric) error
}
