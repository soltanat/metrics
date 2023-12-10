package poller

import (
	"fmt"
	"github.com/soltanat/metrics/internal"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
)

const (
	pollCounterMetricName = "PollCounter"
	randomValueMetricName = "RandomValue"
)

type RuntimePoller struct {
	metrics *sync.Map
}

func NewPoller() *RuntimePoller {
	metrics := &sync.Map{}
	pollCounter := &internal.Metric{
		Name:    pollCounterMetricName,
		Type:    internal.CounterType,
		Counter: 0,
		Gauge:   0,
	}
	randomValue := &internal.Metric{
		Name:    randomValueMetricName,
		Type:    internal.GaugeType,
		Counter: 0,
		Gauge:   0,
	}
	metrics.Store(pollCounterMetricName, pollCounter)
	metrics.Store(randomValueMetricName, randomValue)

	poller := &RuntimePoller{metrics: metrics}

	return poller
}

func (p *RuntimePoller) Get() []internal.Metric {
	var metrics []internal.Metric

	p.metrics.Range(func(key, value interface{}) bool {
		metrics = append(metrics, *value.(*internal.Metric))
		return true
	})

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

		p.metrics.Store(metricName, &internal.Metric{
			Name:    metricName,
			Type:    internal.GaugeType,
			Counter: 0,
			Gauge:   metricValue,
		})
	}

	if m, ok := p.metrics.Load(randomValueMetricName); ok {
		m.(*internal.Metric).SetGauge(rand.Float64())
	}
	if m, ok := p.metrics.Load(pollCounterMetricName); ok {
		m.(*internal.Metric).IncCounter()
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
