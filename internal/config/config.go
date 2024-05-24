package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Address  string         `yaml:"address"`
	LogLevel string         `yaml:"log_level" env-default:"info"`
	Storage  string         `yaml:"storage"`
	DB       DatabaseConfig `yaml:"database"`
}

func MustLoad() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(`config.yaml`, &cfg)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}
