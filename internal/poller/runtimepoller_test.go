package poller

import (
	"testing"
)

func TestNewPoller(t *testing.T) {
	tests := []struct {
		name string
		//want *runtimePoller
	}{
		{
			"new poller",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPoller()
			if got == nil {
				t.Errorf("new poller return nil")
			}
			if _, ok := got.metrics[pollCounterMetricName]; !ok {
				t.Errorf("new poller have not metric %s", pollCounterMetricName)
			}
			if _, ok := got.metrics[randomValueMetricName]; !ok {
				t.Errorf("new poller have not metric %s", randomValueMetricName)
			}
		})
	}
}

func TestPeriodicRuntimePoller_Get(t *testing.T) {
	p := NewPoller()
	err := p.Poll()
	if err != nil {
		t.Errorf("Poll error %v", err)
	}
	got := p.Get()
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
			p := NewPoller()

			previousMetrics := p.metrics

			err := p.Poll()
			if err != nil {
				t.Errorf("Poll error %v", err)
			}

			for key := range gaugeMetrics {
				if _, ok := p.metrics[key]; !ok {
					t.Errorf("poller not polled metric %s", key)
				}
			}
			m, _ := p.metrics[pollCounterMetricName]
			if m.Counter == previousMetrics[pollCounterMetricName].Counter+1 {
				t.Errorf("counter metric not incremented")
			}

			m, _ = p.metrics[randomValueMetricName]
			if m.Gauge != previousMetrics[randomValueMetricName].Gauge {
				t.Errorf("gauge metric value not changed")
			}
		})
	}
}
