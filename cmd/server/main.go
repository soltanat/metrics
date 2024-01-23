package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soltanat/metrics/internal/db"
	"github.com/soltanat/metrics/internal/logger"

	"github.com/soltanat/metrics/internal/filestorage"
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	ctx := context.Background()

	parseFlags()

	l := logger.Get()

	d, err := db.New(ctx, flagDBAddr)
	if err != nil {
		l.Fatal().Err(err)
	}

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

	h := handler.New(fs, d)

	server := handler.SetupRoutes(h)

	go func() {
		err = server.Start(flagAddr)
		if err != nil {
			l.Error().Err(err)
		}
	}()

	gracefulShutdown()

	err = fs.Stop()
	if err != nil {
		l.Error().Err(err)
	}

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
