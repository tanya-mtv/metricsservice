package main

import (
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/server"
)

func main() {

	cfg, err := config.InitServer()
	if err != nil {

		panic("error initialazing config")
	}
	appLogger := logger.NewAppLogger()
	srv := server.NewServer(appLogger, cfg)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
