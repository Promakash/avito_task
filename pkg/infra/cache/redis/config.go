package redis

import "time"

type Config struct {
	Host         string        `env:"CACHE_HOST" yaml:"host" env-required:"true"`
	Port         int           `env:"CACHE_PORT" yaml:"port" env-required:"true"`
	Password     string        `env:"CACHE_PASSWORD" yaml:"password" env-required:"true"`
	TTL          time.Duration `env:"CACHE_TTL" yaml:"TTL" env-default:"30min"`
	WriteTimeout time.Duration `env:"CACHE_WRITE_TIMEOUT" yaml:"write_timeout" env-default:"3s"`
	ReadTimeout  time.Duration `env:"CACHE_READ_TIMEOUT" yaml:"read_timeout" env-default:"2s"`
}
