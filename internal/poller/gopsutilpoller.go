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

type GoPSUtilPoller struct {
	metricsChan chan *model.Metric
}

func NewGoPSUtilPoller(metricsChan chan *model.Metric) *GoPSUtilPoller {
	return &GoPSUtilPoller{metricsChan: metricsChan}
}

func (p *GoPSUtilPoller) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			v, err := mem.VirtualMemory()
			if err != nil {
				return fmt.Errorf("failed to get memory stats: %v", err)
			}

			p.metricsChan <- model.NewGauge(totalMemoryMetricName, float64(v.Total))
			p.metricsChan <- model.NewGauge(freeMemoryMetricName, float64(v.Free))

			c, err := cpu.Percent(0, false)
			if err != nil {
				return fmt.Errorf("failed to get cpu stats: %v", err)
			}
			p.metricsChan <- model.NewGauge(cpuUtilization1MetricName, c[0])
		}
	}
}
