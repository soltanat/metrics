package poller

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPoller(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"new poller",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRuntimePoller()
			assert.NotNil(t, got)
		})
	}
}

func TestPeriodicRuntimePoller_poll(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"check polled metrics"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRuntimePoller()
			assert.NotNil(t, p)

			ctx, cancel := context.WithCancel(context.Background())

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := p.RunPoller(ctx, time.Second)
				assert.NoError(t, err)
			}()

			cancel()
			wg.Wait()
		})
	}
}
