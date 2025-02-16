package config

import (
	"avito_shop/pkg/infra"
	"avito_shop/pkg/infra/cache/redis"
	pkglog "avito_shop/pkg/log"
	"time"
)

type HTTPConfig struct {
	Address      string        `env:"SERVER_ADDRESS" yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" yaml:"idle_timeout" env-default:"30s"`
}

type Config struct {
	HTTPServer HTTPConfig           `yaml:"http_server" env-required:"true"`
	PG         infra.PostgresConfig `yaml:"postgres" env-required:"true"`
	Redis      redis.Config         `yaml:"redis" env-required:"true"`
	Logger     pkglog.Config        `yaml:"logger" env-required:"true"`
	AuthSecret string               `env:"AUTH_SECRET" env-required:"true"`
}

func (c Config) Redact() Config {
	const privateData = "PRIVATE"

	c.PG.User = privateData
	c.PG.Password = privateData
	c.PG.Host = privateData

	c.Redis.Host = privateData
	c.Redis.Password = privateData

	c.AuthSecret = privateData

	return c
}
