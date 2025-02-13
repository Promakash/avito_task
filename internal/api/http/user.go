package http

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	libmiddleware "avito_shop/internal/lib/middleware"
	"avito_shop/internal/usecases"
	"avito_shop/pkg/http/handlers"
	resp "avito_shop/pkg/http/responses"
	pkglog "avito_shop/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	logger  *slog.Logger
	service usecases.User
}

func NewUserHandler(logger *slog.Logger, service usecases.User) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

const getInfoPath = "/info"

func (h *UserHandler) WithSecuredUserHandlers(authService usecases.Auth) handlers.RouterOption {
	return func(r chi.Router) {
		r.With(libmiddleware.WithTokenAuth(authService)).Group(func(r chi.Router) {
			handlers.AddHandler(r.Get, getInfoPath, h.getInfo)
		})
	}
}

func (h *UserHandler) getInfo(r *http.Request) resp.Response {
	const op = "UserHandler.getInfo"
	uid, err := libmiddleware.GetUserIDFromContext(r)
	if err != nil {
		return domain.HandleResult(err, nil)
	}

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.Int("user_id", uid),
	)

	info, err := h.service.GetInfoByID(r.Context(), uid)
	if err != nil {
		log.Error("error while collecting user info: ", pkglog.Err(err))
		return domain.HandleResult(err, nil)
	}

	return domain.HandleResult(nil, types.CreateGetInfoResponse(info))
}
