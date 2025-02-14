package http

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	"avito_shop/internal/usecases"
	"avito_shop/pkg/http/handlers"
	resp "avito_shop/pkg/http/responses"
	pkglog "avito_shop/pkg/log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AuthHandler struct {
	logger  *slog.Logger
	service usecases.Auth
}

func NewAuthHandler(logger *slog.Logger, service usecases.Auth) *AuthHandler {
	return &AuthHandler{
		logger:  logger,
		service: service,
	}
}

const postAuthPath = "/auth"

func (h *AuthHandler) WithAuthHandlers() handlers.RouterOption {
	return func(r chi.Router) {
		handlers.AddHandler(r.Post, postAuthPath, h.postAuth)
	}
}

// @Summary	Аутентификация и получение JWT-токена
// @Accept		json
// @Produce	json
// @Param		body	body		types.PostAuthRequest	true	"Данные пользователя для авторизации"
// @Success	200		{object}	types.PostAuthResponse	"Успешная аутентификация"
// @Failure	400		{object}	responses.ErrorResponse	"Неверный запрос"
// @Failure	401		{object}	responses.ErrorResponse	"Неавторизован"
// @Failure	500		{object}	responses.ErrorResponse	"Внутренняя ошибка сервера"
// @Router		/api/auth [post]
func (h *AuthHandler) postAuth(r *http.Request) resp.Response {
	const op = "AuthHandler.postAuth"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	req, err := types.CreatePostAuthRequest(r)
	if err != nil {
		log.Warn("error while forming request", pkglog.Err(err))
		return domain.HandleResult(domain.ErrBadRequest, nil)
	}

	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Warn("error with user login", pkglog.Err(err))
		return domain.HandleResult(err, nil)
	}

	return domain.HandleResult(nil, types.CreatePostAuthResponse(token))
}
