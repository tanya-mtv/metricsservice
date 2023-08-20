package main

import (
	"github.com/tanya-mtv/metricsservice/internal/agent"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

func main() {
	cfg, err := config.InitAgent()
	if err != nil {

		panic("error initialazing config")
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	ag := agent.NewAgent(cfg, appLogger)

	if err := ag.Run(); err != nil {
		appLogger.Fatal(err)
	}

}
