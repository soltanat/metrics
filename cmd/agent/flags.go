package main

import (
	"encoding/json"
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
var flagConfig string
var flagGRPCServerAddr string

type Config struct {
	Addr           string `env:"ADDRESS" json:"addr"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	Key            string `env:"KEY" json:"key"`
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config         string `env:"CONFIG"`
	GRPCServerAddr string `env:"GRPC_SERVER_ADDR" json:"grpc_server_addr"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.IntVar(&flagReportInterval, "r", 10, "send metrics report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll metrics interval")
	flag.StringVar(&flagKey, "k", "", "key for signature")
	flag.IntVar(&flagRateLimit, "l", 1, "rate limit")
	flag.StringVar(&flagCryptoKey, "crypto-key", "./public_key.pem", "crypto key")
	flag.StringVar(&flagConfig, "config", "", "config path")
	flag.StringVar(&flagGRPCServerAddr, "g", ":9090", "grpc server address")
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
	if cfg.GRPCServerAddr != "" {
		flagGRPCServerAddr = cfg.GRPCServerAddr
	}

	if cfg.Config != "" {
		flagConfig = cfg.Config
	}
	if flagConfig != "" {
		_, err := os.Stat(flagConfig)
		if err != nil {
			l.Fatal().Err(err)
		}

		f, err := os.ReadFile(flagConfig)
		if err != nil {
			l.Fatal().Err(err)
		}
		jsonConfig := Config{}
		err = json.Unmarshal(f, &jsonConfig)
		if err != nil {
			l.Fatal().Err(err)
		}

		if flagAddr == "" && jsonConfig.Addr != "" {
			flagAddr = jsonConfig.Addr
		}
		if flagReportInterval == 0 && jsonConfig.ReportInterval != 0 {
			flagReportInterval = jsonConfig.ReportInterval
		}
		if flagPollInterval == 0 && jsonConfig.PollInterval != 0 {
			flagPollInterval = jsonConfig.PollInterval
		}
		if flagKey == "" && jsonConfig.Key != "" {
			flagKey = jsonConfig.Key
		}
		if flagRateLimit == 0 && jsonConfig.RateLimit != 0 {
			flagRateLimit = jsonConfig.RateLimit
		}
		if flagCryptoKey == "" && jsonConfig.CryptoKey != "" {
			flagCryptoKey = jsonConfig.CryptoKey
		}
		if flagGRPCServerAddr == "" && jsonConfig.GRPCServerAddr != "" {
			flagGRPCServerAddr = jsonConfig.GRPCServerAddr
		}
	}
}
