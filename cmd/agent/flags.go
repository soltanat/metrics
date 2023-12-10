package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	_ "github.com/caarlos0/env/v6"
	"log"
)

var flagAddr string
var flagReportInterval int
var flagPollInterval int

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Addr != "" {
		flagAddr = cfg.Addr
	}
	if cfg.ReportInterval != 0 {
		flagReportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		flagPollInterval = cfg.PollInterval
	}
}
