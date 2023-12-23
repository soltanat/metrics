package main

import (
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	parseFlags()

	l := logger.Get()

	s := storage.NewMemStorage()
	h := handler.New(s)

	err := handler.SetupRoutes(h).Start(flagAddr)
	if err != nil {
		l.Fatal().Err(err)
	}
}
