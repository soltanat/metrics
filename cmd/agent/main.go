package main

import (
	"context"
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

const (
	address               = "http://localhost:8080"
	defaultPollInterval   = time.Second * 2
	defaultReportInterval = time.Second * 10
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
		select {
		case <-ticker.C:
			log.Printf("call poller")
			err := poller.Poll()
			if err != nil {
				cancel()
				return
			}
		case <-ctx.Done():
			log.Printf("poller get context.Done")
			ticker.Stop()
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(reportInterval)
		log.Printf("runned reporter goroutine")
		select {
		case <-ticker.C:
			log.Printf("call reporter")
			err := reporter.Report()
			if err != nil {
				cancel()
				return
			}
		case <-ctx.Done():
			log.Printf("reporter get context.Done")
			ticker.Stop()
			return
		}
	}()

	wg.Wait()
}

func main() {
	pollerInst := poller.NewPoller()
	reporterInst := reporter.New(pollerInst, client.New(address))
	Run(context.Background(), defaultPollInterval, defaultReportInterval, pollerInst, reporterInst)
}
