package config

import (
	"avito_shop/pkg/infra"
	"avito_shop/pkg/infra/cache/redis"
	pkglog "avito_shop/pkg/log"
	"time"
)

type HTTPConfig struct {
	Address      string        `yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type Config struct {
	HTTPServer HTTPConfig           `yaml:"http_server" env-required:"true"`
	PG         infra.PostgresConfig `yaml:"postgres" env-required:"true"`
	Redis      redis.Config         `yaml:"redis" env-required:"true"`
	Logger     pkglog.Config        `yaml:"logger" env-required:"true"`
	AuthSecret string               `env:"AUTH_SECRET" env-required:"true"`
}
