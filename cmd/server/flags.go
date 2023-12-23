package main

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/soltanat/metrics/internal/logger"
)

var flagAddr string

type Config struct {
	Addr string `env:"ADDRESS"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
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
