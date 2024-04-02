// Package reporter
package reporter

import (
	"context"
	"fmt"
	"time"

	"github.com/soltanat/metrics/internal/model"

	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
)

// Reporter
// Реализует интерфейс Reporter
type Reporter struct {
	client      *client.Client
	metricsChan chan *model.Metric
	limitChan   chan struct{}
}

func New(client *client.Client, metricsChan chan *model.Metric, limitChan chan struct{}) *Reporter {
	reporter := &Reporter{
		client:      client,
		metricsChan: metricsChan,
		limitChan:   limitChan,
	}
	return reporter
}

// Run
// Запускает Reporter
// Отправляет метрики в хранилище с помощью клиента, реализует отправку с повторными попытками и рейт лимитом
func (w *Reporter) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	var metrics []model.Metric

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			if len(metrics) == 0 {
				continue
			}
			w.limitChan <- struct{}{}
			err := internal.Backoff(func() error {
				return w.client.Updates(metrics)
			})
			<-w.limitChan
			if err != nil {
				return fmt.Errorf("update metrics error: %w", err)
			}
		case m := <-w.metricsChan:
			metrics = append(metrics, *m)
		}
	}
}
