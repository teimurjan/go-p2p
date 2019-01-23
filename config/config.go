package config

import "github.com/kelseyhightower/envconfig"

// Config is a config variables for the app
type Config struct {
	LogFile  string `envconfig:"LOG_FILE"`
	Port     string `envconfig:"BACKEND_PORT"`
	RedisUrl string `envconfig:"REDIS_URL"`
}

// NewConfig creates new Config instance
func NewConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return &c, err
	}
	return &c, nil
}
