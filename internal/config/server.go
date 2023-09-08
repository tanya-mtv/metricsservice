package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type ConfigServer struct {
	Port     string `env:"ADDRESS"`
	Interval int    `env:"STORE_INTERVAL"`
	FileName string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
	DSN      string `env:"DATABASE_DSN"`
	Logger   *logger.Config
}

func InitServer() (*ConfigServer, error) {

	var flagRunAddr string
	var flagInterval int
	var flagFileName string
	var flagRestore bool
	var flagDSN string

	cfg := &ConfigServer{}
	env.Parse(cfg)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagInterval, "i", 300, "Saved interval")
	flag.StringVar(&flagFileName, "f", "/tmp/metrics-db.json", "storage file")
	flag.BoolVar(&flagRestore, "r", true, "need of sviving")

	flag.StringVar(&flagDSN, "d", "sslmode=disable host=localhost port=5432 dbname = yametrics user=dbuser password=password123", "connection to database")
	// flag.StringVar(&flagDSN, "d", "", "connection to database")

	flag.Parse()

	if cfg.Port == "" {
		cfg.Port = flagRunAddr
	}
	if cfg.FileName == "" {
		cfg.FileName = flagFileName
	}
	if cfg.Interval == 0 {
		cfg.Interval = flagInterval
	}

	if !cfg.Restore {
		cfg.Restore = flagRestore
	}

	if cfg.DSN == "" {
		cfg.DSN = flagDSN
	}

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}

	cfg.Logger = cfglog

	return cfg, nil
}
