package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/poller"

	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/reporter"
)

func Run(
	ctx context.Context, pollInterval, reportInterval time.Duration, pollers []internal.Poll, reporter internal.Reporter,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g := new(errgroup.Group)

	go gracefulStop(ctx, cancel)

	for i := 0; i < len(pollers); i++ {
		i := i
		g.Go(func() error {
			return pollers[i].Run(ctx, pollInterval)
		})
		g.Go(func() error {
			return reporter.Run(ctx, reportInterval)
		})
	}

	err := g.Wait()
	if err != nil {
		l := logger.Get()
		l.Error().Err(err).Msg("run error")
	}
}

func gracefulStop(ctx context.Context, cancelFunc context.CancelFunc) {
	l := logger.Get()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	l.Info().Msg("ran graceful shutdown goroutine")
	select {
	case <-exit:
		l.Info().Msg("graceful shutdown receive signal")
		cancelFunc()
		return
	case <-ctx.Done():
		l.Info().Msg("graceful shutdown got context.Done")
		cancelFunc()
		return
	}
}

func main() {
	parseFlags()

	metricsChan := make(chan *model.Metric)

	runtimePollerInst := poller.NewRuntimePoller(metricsChan)
	goPSUtilPollerInst := poller.NewGoPSUtilPoller(metricsChan)

	pollers := []internal.Poll{runtimePollerInst, goPSUtilPollerInst}

	addr := fmt.Sprintf("http://%s", flagAddr)
	transport := http.DefaultTransport
	transport = &client.GzipTransport{Transport: transport}
	if flagKey != "" {
		transport = &client.SignatureTransport{Transport: transport, Key: flagKey}
	}
	transport = &client.LoggingTransport{Transport: transport}
	cli := client.New(addr, transport)

	reporterInst := reporter.New(cli, metricsChan, make(chan struct{}, flagRateLimit))
	Run(
		context.Background(),
		time.Second*time.Duration(flagPollInterval),
		time.Second*time.Duration(flagReportInterval),
		pollers,
		reporterInst,
	)
}
