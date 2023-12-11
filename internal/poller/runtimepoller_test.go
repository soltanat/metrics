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

			previousMetrics := p.storage

			err = p.Poll()
			assert.NoError(t, err)

			for key := range gaugeMetrics {
				if _, err := p.storage.GetGauge(key); err != nil {
					assert.NoError(t, err)
				}
			}
			if m, err := p.storage.GetCounter(pollCounterMetricName); err == nil {
				if p, err := previousMetrics.GetCounter(pollCounterMetricName); err != nil {
					assert.NoError(t, err)
					if m.Counter == p.Counter+1 {
						t.Errorf("%s metric not incremented", pollCounterMetricName)
					}
				}
			} else {
				t.Errorf("%s metric not exist", pollCounterMetricName)
			}

			if m, err := p.storage.GetGauge(randomValueMetricName); err == nil {
				if p, err := previousMetrics.GetGauge(randomValueMetricName); err == nil {
					if m.Gauge != p.Gauge {
						t.Errorf("%s metric not updated", randomValueMetricName)
					}
				}
			} else {
				t.Errorf("%s not exist", randomValueMetricName)
			}
		})
	}
}
