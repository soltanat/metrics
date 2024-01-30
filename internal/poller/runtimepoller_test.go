package poller

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/soltanat/metrics/internal/model"

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
			got := NewRuntimePoller(make(chan *model.Metric))
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
			metricsChan := make(chan *model.Metric)
			p := NewRuntimePoller(metricsChan)
			assert.NotNil(t, p)

			//prevCounter, err := p.storage.GetCounter(pollCounterMetricName)
			//assert.NoError(t, err)
			//prevRandom, err := p.storage.GetGauge(randomValueMetricName)
			//assert.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := p.Run(ctx, time.Second)
				assert.NoError(t, err)
			}()

			cancel()
			wg.Wait()
		})
	}
}
