package main

import (
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v6"

	"github.com/soltanat/metrics/internal/logger"
)

var flagAddr string
var flagReportInterval int
var flagPollInterval int
var flagKey string
var flagRateLimit int
var flagCryptoKey string

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.StringVar(&flagKey, "k", "", "key for signature")
	flag.IntVar(&flagRateLimit, "l", 1, "rate limit")
	flag.StringVar(&flagCryptoKey, "crypto-key", "./public_key.pem", "crypto key")
	flag.Parse()

	var cfg Config
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			l.Fatal().Err(err)
		}
	}
	err := env.Parse(&cfg)
	if err != nil {
		l.Fatal().Err(err)
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
	if cfg.Key != "" {
		flagKey = cfg.Key
	}
	if cfg.RateLimit != 0 {
		flagRateLimit = cfg.RateLimit
		if flagRateLimit < 1 {
			flagRateLimit = 1
		}
	}
	if cfg.CryptoKey != "" {
		flagCryptoKey = cfg.CryptoKey
	}
}
