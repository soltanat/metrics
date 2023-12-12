package poller

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			got, err := NewPoller()
			assert.NoError(t, err)
			if got == nil {
				t.Errorf("new poller return nil")
				return
			}
			if got.storage == nil {
				t.Errorf("poller metrics storage is nil")
				return
			}
			if _, err := got.storage.GetCounter(pollCounterMetricName); err != nil {
				assert.NoError(t, err)
			}
			if _, err := got.storage.GetGauge(randomValueMetricName); err != nil {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPeriodicRuntimePoller_Get(t *testing.T) {
	p, err := NewPoller()
	assert.NoError(t, err)
	err = p.Poll()
	assert.NoError(t, err)
	got, err := p.Get()
	assert.NoError(t, err)
	if len(got) != len(gaugeMetrics)+2 {
		t.Errorf("partly get metric from poller")
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
			p, err := NewPoller()
			assert.NoError(t, err)

			prevCounter, err := p.storage.GetCounter(pollCounterMetricName)
			assert.NoError(t, err)
			prevRandom, err := p.storage.GetGauge(randomValueMetricName)
			assert.NoError(t, err)

			err = p.Poll()
			assert.NoError(t, err)

			for key := range gaugeMetrics {
				if _, err := p.storage.GetGauge(key); err != nil {
					assert.NoError(t, err)
				}
			}
			m, err := p.storage.GetCounter(pollCounterMetricName)
			assert.NoError(t, err)
			assert.Equal(t, m.Counter, prevCounter.Counter+1)

			m, err = p.storage.GetGauge(randomValueMetricName)
			assert.NoError(t, err)
			assert.NotEqual(t, m.Gauge, prevRandom.Gauge)
		})
	}
}
