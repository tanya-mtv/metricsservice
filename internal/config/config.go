package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ConfigServer struct {
	Port string
}

type ConfigAgent struct {
	Port           string
	ReportInterval int
	PollInterval   int
}

func InitConfigServer() (*ConfigServer, error) {

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

func InitConfigAgent() (*ConfigAgent, error) {
	var flagRunAddr string
	var pollInterval int
	var reportInterval int

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = envRunAddr
	}

	if envreportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		// if envreportInterval := os.Getenv("REPORT_INTERVAL"); envreportInterval != "" {
		envreportIntervalInt, err := strconv.Atoi(envreportInterval)

		if err != nil {
			fmt.Println("Can't parse value reportInterval to Int")
			envreportIntervalInt = reportInterval
		}
		reportInterval = envreportIntervalInt
	}

	if envpollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		// if envpollInterval := os.Getenv("POLL_INTERVAL"); envpollInterval != "" {
		envpollIntervalInt, err := strconv.Atoi(envpollInterval)
		if err != nil {
			fmt.Println("Can't parse value pollInterval to Int")
			pollInterval = envpollIntervalInt
		}

		pollInterval = envpollIntervalInt
	}

	cfg := &ConfigAgent{

		Port:           flagRunAddr,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}

	return cfg, nil
}
