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

	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/poller"
	"github.com/soltanat/metrics/internal/reporter"
)

func Run(
	ctx context.Context, pollInterval, reportInterval time.Duration, poller internal.Poll, reporter internal.Reporter,
) {
	l := logger.Get()

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		l.Info().Msg("ran gs goroutine")
		select {
		case <-exit:
			l.Info().Msg("gs receive signal")
			cancel()
			return
		case <-ctx.Done():
			l.Info().Msg("gs got context.Done")
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(pollInterval)
		l.Info().Msg("ran poller goroutine")
		for {
			select {
			case <-ticker.C:
				l.Info().Msg("call poller")
				err := poller.Poll()
				if err != nil {
					l.Error().Err(err).Msg("poller error")
					cancel()
					return
				}
				l.Info().Msg("polled metrics")
			case <-ctx.Done():
				l.Info().Msg("poller got context.Done")
				ticker.Stop()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(reportInterval)
		l.Info().Msg("ran reporter goroutine")
		for {
			select {
			case <-ticker.C:
				l.Info().Msg("call reporter")
				err := reporter.Report()
				if err != nil {
					l.Error().Err(err).Msg("reporter error")
				} else {
					l.Info().Msg("metrics reported")
				}
			case <-ctx.Done():
				l.Info().Msg("reporter got context.Done")
				ticker.Stop()
				return
			}
		}
	}()

	wg.Wait()
}

func main() {
	l := logger.Get()

	parseFlags()

	pollerInst, err := poller.NewPoller()
	if err != nil {
		l.Fatal().Err(err)
		return
	}

	addr := fmt.Sprintf("http://%s", flagAddr)
	transport := http.DefaultTransport
	transport = &client.GzipTransport{Transport: transport}
	if flagKey != "" {
		transport = &client.SignatureTransport{Transport: transport, Key: flagKey}
	}
	transport = &client.LoggingTransport{Transport: transport}
	cli := client.New(addr, transport)

	reporterInst := reporter.New(pollerInst, cli)
	Run(
		context.Background(),
		time.Second*time.Duration(flagPollInterval),
		time.Second*time.Duration(flagReportInterval),
		pollerInst,
		reporterInst,
	)
}
