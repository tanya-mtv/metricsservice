package main

import (
	"github.com/tanya-mtv/metricsservice/internal/agent"
	"github.com/tanya-mtv/metricsservice/internal/config"
)

func main() {
	cfg, err := config.InitAgent()
	if err != nil {

		panic("error initialazing config")
	}

	ag := agent.NewAgent(cfg)

	if err := ag.Run(); err != nil {
		panic(err)
	}

}
