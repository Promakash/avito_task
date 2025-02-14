package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host     string `env:"DATABASE_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"DATABASE_PORT" yaml:"port" env-required:"true"`
	User     string `env:"DATABASE_USER" yaml:"user" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" yaml:"password" env-required:"true"`
	DBName   string `env:"DATABASE_NAME" yaml:"db_name" env-required:"true"`
}

func NewPostgresPool(cfg PostgresConfig) (*pgxpool.Pool, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	dbPool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("can't create connection to postgres: %w", err)
	}

	return dbPool, nil
}
