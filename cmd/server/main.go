package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soltanat/metrics/internal/filestorage"
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	parseFlags()

	l := logger.Get()

	s := storage.NewMemStorage()

	interval := time.Duration(flagInterval) * time.Second
	fs, err := filestorage.New(s, interval, flagPath, flagRestore)
	if err != nil {
		l.Fatal().Err(err)
	}
	err = fs.Start()
	if err != nil {
		l.Fatal().Err(err)
	}

	h := handler.New(fs)

	server := handler.SetupRoutes(h)

	go func() {
		err = server.Start(flagAddr)
		if err != nil {
			l.Error().Err(err)
		}
	}()

	gracefulShutdown()
	fs.Stop()
	err = server.Close()

	if err != nil {
		l.Error().Err(err)
	}
}

func gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
}
