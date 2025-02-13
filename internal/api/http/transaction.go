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

type TransactionHandler struct {
	logger  *slog.Logger
	service usecases.Transaction
}

func NewTransactionHandler(logger *slog.Logger, service usecases.Transaction) *TransactionHandler {
	return &TransactionHandler{
		logger:  logger,
		service: service,
	}
}

const (
	postSendCoinPath = "/sendCoin"
	getBuyItemPath   = "/buy/{item}"
)

func (h *TransactionHandler) WithSecuredTransactionHandlers(authService usecases.Auth) handlers.RouterOption {
	return func(r chi.Router) {
		r.With(libmiddleware.WithTokenAuth(authService)).Group(func(r chi.Router) {
			handlers.AddHandler(r.Post, postSendCoinPath, h.postSendCoin)
			handlers.AddHandler(r.Get, getBuyItemPath, h.getBuyItem)
		})
	}
}

func (h *TransactionHandler) postSendCoin(r *http.Request) resp.Response {
	const op = "TransactionHandler.postSendCoin"
	uid, err := libmiddleware.GetUserIDFromContext(r)
	if err != nil {
		return domain.HandleResult(err, nil)
	}

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.Int("user_id", uid),
	)

	req, err := types.CreatePostSendCoinRequest(r)
	if err != nil {
		log.Warn("error while forming request", pkglog.Err(err))
		return domain.HandleResult(domain.ErrBadRequest, nil)
	}

	err = h.service.SendCoinByName(r.Context(),
		domain.Transaction{
			From:   uid,
			Amount: req.Amount,
		},
		req.ToUser,
	)
	if err != nil {
		log.Warn("error with sending coins", pkglog.Err(err))
		return domain.HandleResult(err, nil)
	}

	return domain.HandleResult(nil, nil)
}

func (h *TransactionHandler) getBuyItem(r *http.Request) resp.Response {
	const op = "TransactionHandler.getBuyItem"
	uid, err := libmiddleware.GetUserIDFromContext(r)
	if err != nil {
		return domain.HandleResult(err, nil)
	}

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.Int("user_id", uid),
	)

	req, err := types.CreateGetBuyItemRequest(r)
	if err != nil {
		log.Warn("error while forming request", pkglog.Err(err))
		return domain.HandleResult(err, nil)
	}

	err = h.service.BuyItemByName(r.Context(),
		uid,
		req.Item,
	)
	if err != nil {
		log.Warn("error while buying item", pkglog.Err(err))
		return domain.HandleResult(err, nil)
	}

	return domain.HandleResult(nil, nil)
}
