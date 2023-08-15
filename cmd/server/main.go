package main

import (
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/server"
)

func main() {

	cfg, err := config.InitServer()
	if err != nil {

		panic("error initialazing config")
	}

	srv := server.NewServer(cfg)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
