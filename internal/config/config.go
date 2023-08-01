package config

type Config struct {
	// Postgresql *repository.Config      `mapstructure:"postgres"`
	Port string
}

func InitConfig() (*Config, error) {
	// viper.AddConfigPath("configs")
	// viper.SetConfigName("config")

	// if err := viper.ReadInConfig(); err != nil {
	// 	return &Config{}, err
	// }
	cfg := &Config{
		// Postgresql: &repositoryConfig,
		// Port: viper.GetString("port"),
		Port: "8080",
	}

	return cfg, nil
}
