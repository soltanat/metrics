package reporter

import (
	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
)

type Reporter struct {
	Poller internal.Poll
	Client *client.Client
}

func New(poller internal.Poll, client *client.Client) *Reporter {
	reporter := &Reporter{
		Poller: poller,
		Client: client,
	}
	return reporter
}

func (w *Reporter) Report() error {
	metrics, err := w.Poller.Get()
	if err != nil {
		return err
	}
	for _, m := range metrics {
		err := w.Client.Send(&m)
		if err != nil {
			return err
		}
	}
	return nil
}
