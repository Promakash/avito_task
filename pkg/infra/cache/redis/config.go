package redis

import "time"

type Config struct {
	Host         string        `yaml:"host" env-required:"true"`
	Port         int           `yaml:"port" env-required:"true"`
	Password     string        `yaml:"password" env-required:"true"`
	TTL          time.Duration `yaml:"TTL" env-default:"30min"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"3s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"2s"`
}
