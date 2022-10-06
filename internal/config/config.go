package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port             int    `env:"APP_PORT" envDefault:"22800"`
	ConnectionString string `env:"CONNECTION_STRING"`
	RedisPort        int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisHost        string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPass        string `env:"REDIS_PASS" envDefault:""`
}

func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
