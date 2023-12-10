package main

import (
	"context"
	"fmt"
	"github.com/soltanat/metrics/internal"
	"github.com/soltanat/metrics/internal/client"
	"github.com/soltanat/metrics/internal/poller"
	"github.com/soltanat/metrics/internal/reporter"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(
	ctx context.Context, pollInterval, reportInterval time.Duration, poller internal.Poll, reporter internal.Reporter,
) {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		log.Printf("runned gs goroutine")
		select {
		case <-exit:
			log.Printf("gs receive signal")
			cancel()
			return
		case <-ctx.Done():
			log.Printf("gs get context.Done")
			cancel()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(pollInterval)
		log.Printf("runned poller goroutine")
		for {
			select {
			case <-ticker.C:
				log.Printf("call poller")
				err := poller.Poll()
				if err != nil {
					log.Printf("poller error: %v", err)
					cancel()
					return
				}
				log.Printf("polled metrics")
			case <-ctx.Done():
				log.Printf("poller get context.Done")
				ticker.Stop()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(reportInterval)
		log.Printf("runned reporter goroutine")
		for {
			select {
			case <-ticker.C:
				log.Printf("call reporter")
				err := reporter.Report()
				if err != nil {
					log.Printf("reporter error: %v", err)
					cancel()
					return
				}
				log.Printf("metrics reported")
			case <-ctx.Done():
				log.Printf("reporter get context.Done")
				ticker.Stop()
				return
			}
		}
	}()

	wg.Wait()
}

func main() {
	parseFlags()

	pollerInst := poller.NewPoller()
	reporterInst := reporter.New(pollerInst, client.New(fmt.Sprintf("http://%s", flagAddr)))
	Run(context.Background(), flagPollInterval, flagReportInterval, pollerInst, reporterInst)
}
