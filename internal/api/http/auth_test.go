package http

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	"avito_shop/internal/usecases/mocks"
	"avito_shop/pkg/testutils"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPostAuth_Success(t *testing.T) {
	t.Parallel()

	svc := mocks.NewAuth(t)
	h := NewAuthHandler(testutils.NewDummyLogger(), svc)

	req := &types.PostAuthRequest{
		Username: "Avito",
		Password: "12345",
	}
	expectedResp := &types.PostAuthResponse{
		Token: "token",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	svc.On("Login", mock.Anything, req.Username, req.Password).
		Return(expectedResp.Token, nil)

	resp := h.postAuth(httpReq)

	require.Equal(t, http.StatusOK, resp.StatusCode())
	require.Equal(t, expectedResp, resp.GetPayload())
}

func TestPostAuth_BadRequestCases(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(testutils.NewDummyLogger(), nil)

	tests := []struct {
		Name string
		Req  interface{}
	}{
		{"Empty Username", types.PostAuthRequest{Password: "12345"}},
		{"Empty Password", types.PostAuthRequest{Username: "Avito"}},
		{"Empty Request", types.PostAuthRequest{}},
		{"Broken JSON", []byte("{\"username\":\"avito\",\"password\":\"12345\"")},
	}

	for _, test := range tests {
		httpReq := testutils.NewMockJSONRequest(t, test.Req)

		resp := h.postAuth(httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	}
}

func TestPostAuth_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name    string
		Err     error
		ExpCode int
	}{
		{"Unexpected DBError", errors.New("unexpected DBError"), http.StatusInternalServerError},
		{"Hashed password mismatch", domain.ErrUnauthorized, http.StatusUnauthorized},
		{"Token generation error", errors.New("token generation error"), http.StatusInternalServerError},
		{"User exist(on concurrent write could occur)", domain.ErrUserExists, http.StatusUnauthorized},
	}

	for _, test := range tests {
		svc := mocks.NewAuth(t)
		h := NewAuthHandler(testutils.NewDummyLogger(), svc)

		req := types.PostAuthRequest{Username: "Avito", Password: "12345"}
		httpReq := testutils.NewMockJSONRequest(t, req)

		svc.On("Login", mock.Anything, req.Username, req.Password).
			Return("", test.Err)

		resp := h.postAuth(httpReq)

		require.Equal(t, test.ExpCode, resp.StatusCode())
		svc.AssertExpectations(t)
	}
}
