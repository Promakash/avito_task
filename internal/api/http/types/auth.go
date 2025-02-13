package types

import (
	"avito_shop/internal/domain"
	"avito_shop/pkg/http/handlers"
	"errors"
	"fmt"
	"net/http"
)

type PostAuthRequest struct {
	Username domain.UserName `json:"username"`
	Password string          `json:"password"`
}

func CreatePostAuthRequest(r *http.Request) (*PostAuthRequest, error) {
	var req PostAuthRequest
	err := handlers.DecodeRequest(r, &req)
	if err != nil {
		return nil, fmt.Errorf("CreatePostAuthRequest: error while decoding json: %w", err)
	}

	if len(req.Username) == 0 || len(req.Password) == 0 {
		return nil, errors.New("CreatePostAuthRequest: request field is missed")
	}

	return &req, nil
}

type PostAuthResponse struct {
	Token domain.Token `json:"token"`
}

func CreatePostAuthResponse(token domain.Token) *PostAuthResponse {
	return &PostAuthResponse{Token: token}
}
