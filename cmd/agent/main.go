package main

import (
	"github.com/tanya-mtv/metricsservice/internal/agent"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

func main() {
	cfg, err := config.InitConfigAgent()
	if err != nil {

		panic("error initialazing config")
	}

	appLogger := logger.NewAppLogger()
	ag := agent.NewAgent(appLogger, cfg)

	if err := ag.Run(); err != nil {
		panic(err)
	}

}
