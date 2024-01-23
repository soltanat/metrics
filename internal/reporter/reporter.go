package reporter

import (
	"fmt"

	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
)

type Reporter struct {
	poller internal.Poll
	client *client.Client
}

func New(poller internal.Poll, client *client.Client) *Reporter {
	reporter := &Reporter{
		poller: poller,
		client: client,
	}
	return reporter
}

func (w *Reporter) Report() error {
	metrics, err := w.poller.Get()
	if err != nil {
		return fmt.Errorf("get metrics error: %w", err)
	}
	if len(metrics) == 0 {
		return nil
	}
	err = w.client.Updates(metrics)
	if err != nil {
		return fmt.Errorf("update metrics error: %w", err)
	}
	return nil
}
