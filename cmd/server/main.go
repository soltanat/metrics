package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soltanat/metrics/internal/db"
	"github.com/soltanat/metrics/internal/filestorage"
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	ctx := context.Background()

	parseFlags()

	l := logger.Get()

	var s storage.Storage
	var dbConn *pgxpool.Pool

	if flagDBAddr == "" {
		interval := time.Duration(flagInterval) * time.Second
		fs, err := filestorage.New(storage.NewMemStorage(), interval, flagPath)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to create file storage")
		}

		err = fs.Restore(flagRestore)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to restore file storage")
		}

		err = fs.Start()
		if err != nil {
			l.Fatal().Err(err).Msg("unable to start file storage")
		}

		s = fs

		defer func(fs *filestorage.FileStorage) {
			err := fs.Stop()
			if err != nil {
				l.Error().Err(err).Msg("unable to stop file storage")
			}
		}(fs)
	} else {
		err := db.ApplyMigrations(flagDBAddr)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to apply migrations")
		}

		dbConn, err = db.New(ctx, flagDBAddr)
		if err != nil {
			l.Fatal().Err(err).Msg("unable to connect to database")
		}

		s = storage.NewPostgresStorage(dbConn)
		s = storage.NewBackoffPostgresStorage(s)

		defer dbConn.Close()
	}

	h := handler.New(s, dbConn)

	server := handler.SetupRoutes(h, flagKey)

	go func() {
		err := server.Start(flagAddr)
		if err != nil {
			l.Error().Err(err).Msg("unable to start server")
		}
	}()

	go func() {
		err := http.ListenAndServe(flagPprofAddr, nil)
		if err != nil {
			l.Error().Err(err).Msg("unable to listen and serve")
		}
	}()

	gracefulShutdown()

	err := server.Close()
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
