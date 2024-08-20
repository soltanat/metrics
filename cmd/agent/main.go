package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func Run(
	ctx context.Context,
	pollInterval, reportInterval time.Duration,
	pollers []internal.Poll,
	reporter internal.Reporter,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go gracefulStop(ctx, cancel)

	g, ctx := errgroup.WithContext(ctx)

	chs := make([]chan *model.Metric, len(pollers))
	for i := 0; i < len(pollers); i++ {
		chs = append(chs, pollers[i].GetChannel())
	}
	mergedCh := merge(chs...)

	for i := 0; i < len(pollers); i++ {
		i := i
		g.Go(func() error {
			err := pollers[i].RunPoller(ctx, pollInterval)
			return err
		})
		g.Go(func() error {
			err := reporter.RunReporter(ctx, reportInterval, mergedCh)
			return err
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
	l := logger.Get()
	l.Printf("Build version: %s\n", buildVersion)
	l.Printf("Build date: %s\n", buildDate)
	l.Printf("Build commit: %s\n", buildCommit)

	parseFlags()

	runtimePollerInst := poller.NewRuntimePoller()
	goPSUtilPollerInst := poller.NewGoPSUtilPoller()

	pollers := []internal.Poll{runtimePollerInst, goPSUtilPollerInst}

	addr := fmt.Sprintf("http://%s", flagAddr)

	transport := http.DefaultTransport

	transport = &client.GzipTransport{Transport: transport}

	if flagKey != "" {
		transport = &client.SignatureTransport{Transport: transport, Key: flagKey}
	}
	transport = &client.LoggingTransport{Transport: transport}

	if flagCryptoKey != "" {
		key, err := os.ReadFile(flagCryptoKey)
		if err != nil {
			l.Error().Msg("unable to read crypto key")
			return
		}

		transport, err = client.NewRSAEncryptionTransport(transport, key)
		if err != nil {
			l.Error().Err(err).Msg("unable to create crypto transport")
			return
		}
	}

	cli := client.New(addr, transport)

	reporterInst := reporter.New(cli, make(chan struct{}, flagRateLimit))
	Run(
		context.Background(),
		time.Second*time.Duration(flagPollInterval),
		time.Second*time.Duration(flagReportInterval),
		pollers,
		reporterInst,
	)
}

func merge(cs ...chan *model.Metric) chan *model.Metric {
	var wg sync.WaitGroup
	out := make(chan *model.Metric)

	output := func(c <-chan *model.Metric) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
