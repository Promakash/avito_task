package testutils

import (
	"avito_shop/internal/domain"
	libmiddleware "avito_shop/internal/lib/middleware"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func NewMockJSONRequest(t *testing.T, payload interface{}) *http.Request {
	bodyBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest("", "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	return req
}

const GetBuyItemRequestQueryParam = "item"

func NewMockRequestWithItemQueryVal(itemName string) *http.Request {
	req := httptest.NewRequest("", "/", nil)
	chiCtx := chi.NewRouteContext()

	chiCtx.URLParams.Add(GetBuyItemRequestQueryParam, itemName)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	return req
}

func AddUserIDToRequestContext(r *http.Request, id domain.UserID) *http.Request {
	ctx := r.Context()

	newCtx := context.WithValue(ctx, libmiddleware.AuthContextKey, id)

	return r.WithContext(newCtx)
}

func NewMockRequest() *http.Request {
	return httptest.NewRequest("", "/", nil)
}

func SendRequest(t *testing.T, path, method, authToken string, payload interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(context.Background(), method, path, bytes.NewReader(bodyBytes))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	if len(authToken) > 0 {
		req.Header.Set("Authorization", authToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
