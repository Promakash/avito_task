package main

import (
	httpapp "avito_shop/internal/app/http"
	"avito_shop/internal/config"
	"avito_shop/internal/repository/postgres"
	"avito_shop/internal/usecases/service"
	pkgconfig "avito_shop/pkg/config"
	"avito_shop/pkg/infra"
	pkgredis "avito_shop/pkg/infra/cache/redis"
	pkglog "avito_shop/pkg/log"
	"avito_shop/pkg/shutdown"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

const (
	ConfigEnvVar = "SHOP_CONFIG"
	APIPath      = "/api"
)

func main() {
	cfg := config.Config{}
	pkgconfig.MustLoad(ConfigEnvVar, &cfg)

	log, file := pkglog.NewLogger(cfg.Logger)
	defer func() { _ = file.Close() }()
	slog.SetDefault(log)
	log.Info("Starting Avito Shop", slog.Any("config", cfg))

	dbPool, err := infra.NewPostgresPool(cfg.PG)
	if err != nil {
		pkglog.Fatal(log, "error while setting new postgres connection: ", err)
	}
	defer dbPool.Close()

	redisClient, err := pkgredis.NewRedisClient(cfg.Redis)
	if err != nil {
		pkglog.Fatal(log, "error while setting new redis connection: ", err)
	}
	defer pkgredis.ShutdownClient(redisClient)

	userCache := pkgredis.NewRedisService(redisClient, log)

	txRepo := postgres.NewTransactionRepository(dbPool)
	userRepo := postgres.NewUserRepository(dbPool, userCache, cfg.Redis.TTL, cfg.Redis.WriteTimeout)
	merchRepo := postgres.NewMerchRepository(dbPool)

	authService := service.NewAuth(userRepo, cfg.AuthSecret)
	userService := service.NewUser(userRepo)
	txService := service.NewTransaction(txRepo, userRepo, merchRepo)

	httpApp := httpapp.New(
		log,
		APIPath,
		authService,
		userService,
		txService,
		cfg.HTTPServer,
	)

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return shutdown.ListenSignal(ctx, log)
	})

	g.Go(func() error {
		return httpApp.Run()
	})

	g.Go(func() error {
		<-ctx.Done()
		log.Info("Shutdown signal received, stopping server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpApp.Stop(shutdownCtx)
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, shutdown.ErrOSSignal) {
		log.Error("Exit reason", slog.String("error", err.Error()))
	}
}
