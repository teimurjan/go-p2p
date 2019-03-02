package config

import "github.com/kelseyhightower/envconfig"

// Config is a config variables for the app
type Config struct {
	LogFile       string `envconfig:"LOG_FILE"`
	TCPPort       string `envconfig:"BACKEND_TCP_PORT"`
	UDPPort       string `envconfig:"BACKEND_UDP_PORT"`
	HTTPPort      string `envconfig:"BACKEND_HTTP_PORT"`
	RedisURL      string `envconfig:"REDIS_URL"`
	FileSourceDir string `envconfig:"FILE_DIR_PATH"`
}

// NewConfig creates new Config instance
func NewConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return &c, err
	}
	return &c, nil
}
