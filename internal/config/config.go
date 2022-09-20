package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port             int    `env:"PORT" envDefault:"22800"`
	ConnectionString string `env:"CONNECTION_STRING"`
}

func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
