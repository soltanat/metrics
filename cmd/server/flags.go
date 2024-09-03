package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/soltanat/metrics/internal/logger"
)

var (
	flagAddr           string
	flagPprofAddr      string
	flagInterval       int
	flagPath           string
	flagRestore        bool
	flagDBAddr         string
	flagKey            string
	flagCryptoKey      string
	flagConfig         string
	flagTrustedSubnet  string
	flagGRPCServerAddr string
)

type Config struct {
	Addr           string `env:"ADDRESS" json:"addr"`
	Interval       int    `env:"STORE_INTERVAL" json:"interval"`
	Path           string `env:"FILE_STORAGE_PATH" json:"path"`
	Restore        bool   `env:"RESTORE" json:"restore"`
	DBAddr         string `env:"DATABASE_DSN" json:"db_addr"`
	Key            string `env:"KEY" json:"key"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	Config         string `env:"CONFIG"`
	TrustedSubnet  string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	GRPCServerAddr string `env:"GRPC_SERVER_ADDR" json:"grpc_server_addr"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port metrics http server")
	flag.StringVar(&flagPprofAddr, "p", "localhost:6060", "address and port pprof http server")
	flag.IntVar(&flagInterval, "i", 300, "store metrics interval")
	flag.StringVar(&flagPath, "f", "/tmp/metrics-db.json", "path to store metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore metrics from file")
	flag.StringVar(&flagDBAddr, "d", "", "database dsn")
	flag.StringVar(&flagKey, "k", "", "key for signature")
	flag.StringVar(&flagCryptoKey, "crypto-key", "./private_key.pem", "crypto key")
	flag.StringVar(&flagConfig, "config", "", "config path")
	flag.StringVar(&flagTrustedSubnet, "t", "", "trusted subnet")
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
	if cfg.DBAddr != "" {
		flagDBAddr = cfg.DBAddr
	}
	if cfg.Key != "" {
		flagKey = cfg.Key
	}
	if cfg.CryptoKey != "" {
		flagCryptoKey = cfg.CryptoKey
	}
	if cfg.TrustedSubnet != "" {
		flagTrustedSubnet = cfg.TrustedSubnet
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
		if flagDBAddr == "" && jsonConfig.DBAddr != "" {
			flagDBAddr = jsonConfig.DBAddr
		}
		if flagKey == "" && jsonConfig.Key != "" {
			flagKey = jsonConfig.Key
		}
		if flagCryptoKey == "" && jsonConfig.CryptoKey != "" {
			flagCryptoKey = jsonConfig.CryptoKey
		}
		if flagTrustedSubnet == "" && jsonConfig.TrustedSubnet != "" {
			flagTrustedSubnet = jsonConfig.TrustedSubnet
		}
		if flagGRPCServerAddr == "" && jsonConfig.GRPCServerAddr != "" {
			flagGRPCServerAddr = jsonConfig.GRPCServerAddr
		}
	}
}
