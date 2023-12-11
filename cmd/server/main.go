package main

import (
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	parseFlags()

	s := storage.NewMemStorage()
	h := handler.New(s)

	e := handler.SetupRoutes(h)
	e.Logger.Fatal(e.Start(flagAddr))
}
