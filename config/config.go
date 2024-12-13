package config

import (
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func New() *Config {
	var Cfg Config
	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}
	return &Cfg
}
