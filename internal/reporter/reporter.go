// Package reporter
package reporter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/soltanat/metrics/internal/model"

	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
)

// Reporter
// Реализует интерфейс Reporter
type Reporter struct {
	client    *client.Client
	limitChan chan struct{}
}

func New(client *client.Client, limitChan chan struct{}) *Reporter {
	reporter := &Reporter{
		client:    client,
		limitChan: limitChan,
	}
	return reporter
}

// Run
// Запускает Reporter
// Отправляет метрики в хранилище с помощью клиента, реализует отправку с повторными попытками и рейт лимитом
func (w *Reporter) RunReporter(ctx context.Context, interval time.Duration, ch chan *model.Metric) error {
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
				for i := 0; i < len(metrics); i++ {
					chunk := 10
					if i+chunk > len(metrics) {
						chunk = len(metrics) - i
					}
					err := w.client.Updates(metrics[i : i+chunk])
					if err != nil {
						if errors.Is(err, client.ErrForbidden) {
							return model.ErrForbidden
						}
						return err
					}
					i += chunk
				}
				return nil
			}, model.ErrForbidden)
			<-w.limitChan
			if err != nil {
				return fmt.Errorf("update metrics error: %w", err)
			}
			metrics = metrics[:0]
		case m := <-ch:
			metrics = append(metrics, *m)
		}
	}
}
