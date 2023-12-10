package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var flagAddr string

type Config struct {
	Addr string `env:"ADDRESS"`
}

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Addr != "" {
		flagAddr = cfg.Addr
	}
}
