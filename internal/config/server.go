package config

import (
	"flag"
	"fmt"
	"os"
)

type ConfigServer struct {
	Port string
}

func InitServer() (*ConfigServer, error) {

	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = envRunAddr
	}

	fmt.Println("flagRunAddr", flagRunAddr)
	cfg := &ConfigServer{
		Port: flagRunAddr,
	}

	return cfg, nil
}
