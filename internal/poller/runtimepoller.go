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

// RuntimePoller
// Реализует интерфейс Poll для сбора рантайм метрик
type RuntimePoller struct {
	metricsChan chan *model.Metric
}

func NewRuntimePoller() *RuntimePoller {
	metricsChan := make(chan *model.Metric)
	return &RuntimePoller{metricsChan: metricsChan}
}

// Run
// Запускает сбор метрик
// interval - интервал сбора метрик
func (p *RuntimePoller) RunPoller(ctx context.Context, interval time.Duration) error {
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

			metrics := make([]*model.Metric, 0)

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

				metrics = append(metrics, model.NewGauge(metricName, metricValue))
			}

			metrics = append(metrics, model.NewGauge(randomValueMetricName, rand.Float64()))
			metrics = append(metrics, model.NewGauge(pollCounterMetricName, 1))

			if err := p.sendMetric(ctx, metrics); err != nil {
				return err
			}
		}
	}
}

func (p *RuntimePoller) sendMetric(ctx context.Context, metric []*model.Metric) error {
	for i := 0; i < len(metric); i++ {
		select {
		case p.metricsChan <- metric[i]:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (p *RuntimePoller) GetChannel() chan *model.Metric {
	return p.metricsChan
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
