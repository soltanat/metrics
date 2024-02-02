package internal

import (
	"context"
	"time"
)

type Poll interface {
	Run(ctx context.Context, interval time.Duration) error
}

type Reporter interface {
	Run(ctx context.Context, interval time.Duration) error
}
