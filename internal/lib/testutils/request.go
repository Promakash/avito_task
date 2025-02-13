package testutils

import (
	"avito_shop/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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

	newCtx := context.WithValue(ctx, "user_id", id)

	return r.WithContext(newCtx)
}

func NewMockRequest() *http.Request {
	return httptest.NewRequest("", "/", nil)
}
