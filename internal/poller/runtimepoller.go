package poller

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
	"math/rand"
	"reflect"
	"runtime"
)

const (
	pollCounterMetricName = "PollCounter"
	randomValueMetricName = "RandomValue"
)

type RuntimePoller struct {
	metrics map[string]internal.Metric
}

func NewPoller() *RuntimePoller {
	metrics := make(map[string]internal.Metric)
	pollCounter := internal.Metric{
		Name:    pollCounterMetricName,
		Type:    internal.CounterType,
		Counter: 0,
		Gauge:   0,
	}
	randomValue := internal.Metric{
		Name:    randomValueMetricName,
		Type:    internal.GaugeType,
		Counter: 0,
		Gauge:   0,
	}
	metrics[pollCounterMetricName] = pollCounter
	metrics[randomValueMetricName] = randomValue

	poller := &RuntimePoller{metrics: metrics}

	return poller
}

func (p *RuntimePoller) Get() []internal.Metric {
	var metrics []internal.Metric
	for _, value := range p.metrics {
		metrics = append(metrics, value)
	}
	return metrics
}

func (p *RuntimePoller) Poll() error {
	runtimeMetrics := &runtime.MemStats{}
	runtime.ReadMemStats(runtimeMetrics)

	v := reflect.ValueOf(*runtimeMetrics)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		metricName := t.Field(i).Name
		if _, ok := gaugeMetrics[metricName]; !ok {
			continue
		}

		var metricValue float64
		switch v.Field(i).Interface().(type) {
		case uint64:
			metricValue = float64(v.Field(i).Interface().(uint64))
		case uint32:
			metricValue = float64(v.Field(i).Interface().(uint32))
		case float64:
			metricValue = v.Field(i).Interface().(float64)
		default:
			return fmt.Errorf("unkonwn metric type %T", v.Field(i).Interface())
		}

		p.metrics[metricName] = internal.Metric{
			Name:    metricName,
			Type:    internal.GaugeType,
			Counter: 0,
			Gauge:   metricValue,
		}
	}

	if m, ok := p.metrics[randomValueMetricName]; ok {
		m.SetGauge(rand.Float64())
	}
	if m, ok := p.metrics[pollCounterMetricName]; ok {
		m.IncCount()
	}
	return nil
}

var gaugeMetrics = map[string]struct{}{
	"Alloc":         {},
	"BuckHashSys":   {},
	"Frees":         {},
	"GCCPUFraction": {},
	"GCSys":         {},
	"HeapAlloc":     {},
	"HeapIdle":      {},
	"HeapInuse":     {},
	"HeapObjects":   {},
	"HeapReleased":  {},
	"HeapSys":       {},
	"LastGC":        {},
	"Lookups":       {},
	"MCacheInuse":   {},
	"MCacheSys":     {},
	"MSpanInuse":    {},
	"MSpanSys":      {},
	"Mallocs":       {},
	"NextGC":        {},
	"NumForcedGC":   {},
	"NumGC":         {},
	"OtherSys":      {},
	"PauseTotalNs":  {},
	"StackInuse":    {},
	"StackSys":      {},
	"Sys":           {},
	"TotalAlloc":    {},
}
