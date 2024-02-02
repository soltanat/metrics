package poller

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/soltanat/metrics/internal/model"
)

const (
	pollCounterMetricName = "PollCount"
	randomValueMetricName = "RandomValue"
)

type RuntimePoller struct {
	metricsChan chan *model.Metric
}

func NewRuntimePoller(metricsChan chan *model.Metric) *RuntimePoller {
	return &RuntimePoller{metricsChan: metricsChan}
}

func (p *RuntimePoller) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
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

				p.metricsChan <- model.NewGauge(metricName, metricValue)
			}

			p.metricsChan <- model.NewGauge(randomValueMetricName, rand.Float64())
			p.metricsChan <- model.NewCounter(pollCounterMetricName, 1)
		}
	}
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
