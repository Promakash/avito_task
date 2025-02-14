package types

import (
	"avito_shop/internal/domain"
	"avito_shop/pkg/http/handlers"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PostSendCoinRequest struct {
	ToUser domain.UserName `json:"toUser"`
	Amount int             `json:"amount"`
}

func CreatePostSendCoinRequest(r *http.Request) (*PostSendCoinRequest, error) {
	var req PostSendCoinRequest
	err := handlers.DecodeRequest(r, &req)
	if err != nil {
		return nil, fmt.Errorf("PostSendCoinRequest: error while decoding json: %w", err)
	}

	if len(req.ToUser) == 0 || req.Amount == 0 {
		return nil, errors.New("PostSendCoinRequest: invalid request field")
	}

	return &req, nil
}

type GetBuyItemRequest struct {
	Item string
}

func CreateGetBuyItemRequest(r *http.Request) (*GetBuyItemRequest, error) {
	const queryParamName = "item"
	itemName := chi.URLParam(r, queryParamName)
	if itemName == "" {
		return nil, fmt.Errorf("CreateGetBuyItemRequest: invalid query provided: %w", domain.ErrBadRequest)
	}

	return &GetBuyItemRequest{Item: itemName}, nil
}
