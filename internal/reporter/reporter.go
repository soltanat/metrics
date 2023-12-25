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
	for _, m := range metrics {
		err := w.client.Update(&m)
		if err != nil {
			return fmt.Errorf("update metrics error: %w", err)
		}
	}
	return nil
}
