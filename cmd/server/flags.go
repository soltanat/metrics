package main

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/soltanat/metrics/internal/logger"
)

var flagAddr string
var flagInterval int
var flagPath string
var flagRestore bool

type Config struct {
	Addr     string `env:"ADDRESS"`
	Interval int    `env:"STORE_INTERVAL"`
	Path     string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.IntVar(&flagInterval, "i", 300, "store metrics interval")
	flag.StringVar(&flagPath, "f", "/tmp/metrics-db.json", "path to store metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore metrics from file")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		l.Fatal().Err(err)
	}

	if cfg.Addr != "" {
		flagAddr = cfg.Addr
	}
}
