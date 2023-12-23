package main

import (
	"context"
	"fmt"
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
		l.Debug().Msg("ran gs goroutine")
		select {
		case <-exit:
			l.Debug().Msg("gs receive signal")
			cancel()
			return
		case <-ctx.Done():
			l.Debug().Msg("gs got context.Done")
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(pollInterval)
		l.Debug().Msg("ran poller goroutine")
		for {
			select {
			case <-ticker.C:
				l.Debug().Msg("call poller")
				err := poller.Poll()
				if err != nil {
					l.Error().Err(err).Msg("poller error")
					cancel()
					return
				}
				l.Debug().Msg("polled metrics")
			case <-ctx.Done():
				l.Debug().Msg("poller got context.Done")
				ticker.Stop()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(reportInterval)
		l.Debug().Msg("ran reporter goroutine")
		for {
			select {
			case <-ticker.C:
				l.Debug().Msg("call reporter")
				err := reporter.Report()
				if err != nil {
					l.Error().Err(err).Msg("reporter error")
					cancel()
					return
				}
				l.Debug().Msg("metrics reported")
			case <-ctx.Done():
				l.Debug().Msg("reporter got context.Done")
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
	reporterInst := reporter.New(pollerInst, client.New(fmt.Sprintf("http://%s", flagAddr)))
	Run(
		context.Background(),
		time.Second*time.Duration(flagPollInterval),
		time.Second*time.Duration(flagReportInterval),
		pollerInst,
		reporterInst,
	)
}
