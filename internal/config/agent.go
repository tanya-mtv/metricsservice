package config

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type ConfigAgent struct {
	Port           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Logger         *logger.Config
}

func InitAgent() (*ConfigAgent, error) {
	var flagRunAddr string
	var pollInterval int
	var reportInterval int

	cfg := &ConfigAgent{}
	env.Parse(cfg)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	if cfg.PollInterval == 0 {
		cfg.PollInterval = pollInterval
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportInterval
	}

	if cfg.Port == "" {
		cfg.Port = flagRunAddr
	}

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}

	cfg.Logger = cfglog

	return cfg, nil
}
