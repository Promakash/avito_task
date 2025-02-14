package http

import (
	apihttp "avito_shop/internal/api/http"
	"avito_shop/internal/config"
	"avito_shop/internal/usecases"
	"avito_shop/pkg/http/handlers"
	"context"
	"log/slog"
	"net/http"
)

type App struct {
	log    *slog.Logger
	server *http.Server
}

func New(
	log *slog.Logger,
	apiPath string,
	authService usecases.Auth,
	userService usecases.User,
	txService usecases.Transaction,
	cfg config.HTTPConfig,
) *App {
	authHandler := apihttp.NewAuthHandler(
		log,
		authService,
	)

	userHandler := apihttp.NewUserHandler(
		log,
		userService,
	)

	txHandler := apihttp.NewTransactionHandler(
		log,
		txService,
	)

	publicHandler := handlers.NewHandler(
		apiPath,
		handlers.WithRequestID(),
		handlers.WithRecover(),
		handlers.WithLogging(log),
		handlers.WithProfilerHandlers(),
		handlers.WithHealthHandler(),
		handlers.WithSwagger(),
		userHandler.WithSecuredUserHandlers(authService),
		txHandler.WithSecuredTransactionHandlers(authService),
		authHandler.WithAuthHandlers(),
	)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      publicHandler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		log:    log,
		server: srv,
	}
}

func (a *App) Run() error {
	const op = "http.App"

	log := a.log.With(
		slog.String("op", op),
		slog.String("address", a.server.Addr),
	)

	log.Info("HTTP server starting")
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	const op = "http.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("HTTP server shutting down", slog.String("addr", a.server.Addr))
	return a.server.Shutdown(ctx)
}
