package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/soltanat/metrics/internal/model"
)

const (
	totalMemoryMetricName     = "TotalMemory"
	freeMemoryMetricName      = "FreeMemory"
	cpuUtilization1MetricName = "CPUtilization1"
)

// GoPSUtilPoller
// Реализует интерфейс Poll для сбора gopsutil метрик
type GoPSUtilPoller struct {
	metricsChan chan *model.Metric
}

func NewGoPSUtilPoller() *GoPSUtilPoller {
	metricsChan := make(chan *model.Metric)
	return &GoPSUtilPoller{metricsChan: metricsChan}
}

// Run
// Запускает сбор метрик
// Передает собранные метрики в канал metricsChan
// interval - интервал сбора метрик
func (p *GoPSUtilPoller) RunPoller(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			close(p.metricsChan)
			ticker.Stop()
			return nil
		case <-ticker.C:
			v, err := mem.VirtualMemory()
			if err != nil {
				return fmt.Errorf("failed to get memory stats: %v", err)
			}

			metrics := make([]*model.Metric, 0, 3)
			metrics = append(metrics, model.NewGauge(totalMemoryMetricName, float64(v.Total)))
			metrics = append(metrics, model.NewGauge(freeMemoryMetricName, float64(v.Free)))

			c, err := cpu.Percent(0, false)
			if err != nil {
				return fmt.Errorf("failed to get cpu stats: %v", err)
			}

			metrics = append(metrics, model.NewGauge(cpuUtilization1MetricName, c[0]))

			err = p.sendMetric(ctx, metrics)
			if err != nil {
				return fmt.Errorf("failed to send metrics: %v", err)
			}
		}
	}
}

func (p *GoPSUtilPoller) sendMetric(ctx context.Context, metric []*model.Metric) error {
	for i := 0; i < len(metric); i++ {
		select {
		case p.metricsChan <- metric[i]:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (p *GoPSUtilPoller) GetChannel() chan *model.Metric {
	return p.metricsChan
}
