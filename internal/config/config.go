package config

import (
	"flag"
)

type ConfigServer struct {
	// Postgresql *repository.Config      `mapstructure:"postgres"`
	Port string
}

type ConfigAgent struct {
	// Postgresql *repository.Config      `mapstructure:"postgres"`
	Port           string
	ReportInterval int
	PollInterval   int
}

func InitConfigServer() (*ConfigServer, error) {
	// viper.AddConfigPath("configs")
	// viper.SetConfigName("config")

	// if err := viper.ReadInConfig(); err != nil {
	// 	return &Config{}, err
	// }
	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "http://localhost:8080", "address and port to run server")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	cfg := &ConfigServer{
		// Postgresql: &repositoryConfig,
		// Port: viper.GetString("port"),
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

	cfg := &ConfigAgent{
		// Postgresql: &repositoryConfig,
		// Port: viper.GetString("port"),
		Port:           flagRunAddr,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	return cfg, nil
}
