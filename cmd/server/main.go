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

	err := db.ApplyMigrations(flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to apply migrations")
	}

	d, err := db.New(ctx, flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to connect to database")
	}

	s := storage.NewPostgresStorage(d)

	interval := time.Duration(flagInterval) * time.Second
	fs, err := filestorage.New(s, interval, flagPath, flagRestore)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to create file storage")
	}
	err = fs.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to start file storage")
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
		l.Error().Err(err).Msg("unable to stop file storage")
	}

	err = server.Close()
	if err != nil {
		l.Error().Err(err).Msg("unable to close server")
	}
}

func gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
}
