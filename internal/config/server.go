package config

import (
	"flag"
	"os"

	"github.com/tanya-mtv/metricsservice/internal/constants"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type ConfigServer struct {
	Port   string
	Logger *logger.Config
}

func InitServer() (*ConfigServer, error) {

	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = envRunAddr
	}

	cfglog := &logger.Config{
		LogLevel: constants.LogLevel,
		DevMode:  constants.DevMode,
		Type:     constants.Type,
	}

	cfg := &ConfigServer{
		Port:   flagRunAddr,
		Logger: cfglog,
	}

	return cfg, nil
}
