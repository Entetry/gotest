package config

import (
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	Port                   int           `env:"APP_PORT" envDefault:"22800"`
	ConnectionString       string        `env:"CONNECTION_STRING"`
	IsMongo                bool          `env:"IS_MONGO" envDefault:"false"`
	AccessTokenKey         string        `env:"ACCESS_TOKEN_KEY" envDefault:"my-access-token-key"`
	AccessTokenExpiration  time.Duration `env:"ACCESS_TOKEN_EXPIRATION" envDefault:"30m"`
	RefreshTokenKey        string        `env:"REFRESH_TOKEN_KEY" envDefault:"my-refresh-token-key"`
	RefreshTokenExpiration time.Duration `env:"REFRESH_TOKEN_EXPIRATION" envDefault:"228h"`
	Salt                   string        `env:"SALT" envDefault:"loremipsum228"`
}

func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
