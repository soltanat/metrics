package poller

import (
	"fmt"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
	"math/rand"
	"reflect"
	"runtime"
)

const (
	pollCounterMetricName = "PollCounter"
	randomValueMetricName = "RandomValue"
)

type RuntimePoller struct {
	storage storage.Storage
}

func NewPoller() (*RuntimePoller, error) {
	s := storage.NewMemStorage()

	pollCounter := model.NewCounter(pollCounterMetricName, 0)
	randomValue := model.NewGauge(randomValueMetricName, 0)

	err := s.Store(pollCounter)
	if err != nil {
		return nil, err
	}
	err = s.Store(randomValue)
	if err != nil {
		return nil, err
	}

	poller := &RuntimePoller{storage: s}

	return poller, nil
}

func (p *RuntimePoller) Get() ([]model.Metric, error) {
	return p.storage.GetList()
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

		err := p.storage.Store(model.NewGauge(metricName, metricValue))
		if err != nil {
			return err
		}
	}

	if m, err := p.storage.GetGauge(randomValueMetricName); err != nil {
		return err
	} else {
		m.SetGauge(rand.Float64())
	}

	if m, err := p.storage.GetCounter(pollCounterMetricName); err != nil {
		return err
	} else {
		m.IncCounter()
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
