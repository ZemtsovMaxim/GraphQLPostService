package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		DSN  string `yaml:"dsn"`
		Type string `yaml:"type" env-default:"memory"` // postgres или memory
	} `yaml:"database"`
	LogLevel string `yaml:"logLevel" env-default:"info"`
}

// MustLoad загружает конфигурацию из файла config.yaml
func MustLoad() *Config {
	const configPath = "config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
