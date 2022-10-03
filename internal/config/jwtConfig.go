package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type JwtConfig struct {
	AccessTokenKey         string        `env:"ACCESS_TOKEN_KEY" envDefault:"my-access-token-key"`
	AccessTokenExpiration  time.Duration `env:"ACCESS_TOKEN_EXPIRATION" envDefault:"30m"`
	RefreshTokenExpiration time.Duration `env:"REFRESH_TOKEN_EXPIRATION" envDefault:"3000m"`
}

func NewJwtConfig() (*JwtConfig, error) {
	cfg := new(JwtConfig)
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
