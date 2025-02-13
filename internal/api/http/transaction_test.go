package http

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	"avito_shop/internal/lib/testutils"
	"avito_shop/internal/usecases/mocks"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestPostSendCoin_Success(t *testing.T) {
	t.Parallel()

	svc := mocks.NewTransaction(t)
	h := NewTransactionHandler(testutils.NewDummyLogger(), svc)

	uID := 2
	req := types.PostSendCoinRequest{
		ToUser: "Avito",
		Amount: 100,
	}

	httpReq := testutils.NewMockJSONRequest(t, req)
	httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

	svc.On(
		"SendCoinByName",
		mock.Anything,
		domain.Transaction{From: uID, Amount: req.Amount},
		req.ToUser).Return(nil)

	resp := h.postSendCoin(httpReq)

	require.Equal(t, http.StatusOK, resp.StatusCode())
	svc.AssertExpectations(t)
}

func TestPostSendCoin_EmptyContextVal(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(testutils.NewDummyLogger(), nil)

	httpReq := testutils.NewMockRequest()

	resp := h.postSendCoin(httpReq)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
}

func TestPostSendCoin_BadRequestCases(t *testing.T) {
	t.Parallel()

	uID := 2
	tests := []struct {
		Name string
		Req  interface{}
	}{
		{"Empty toUser", types.PostSendCoinRequest{Amount: 100}},
		{"Empty amount", types.PostSendCoinRequest{ToUser: "Avito"}},
		{"Empty Request", types.PostAuthRequest{}},
		{"Broken JSON", []byte("{\"toUser\":\"avito\",\"amount\":\"100\"")},
	}

	for _, test := range tests {
		h := NewTransactionHandler(testutils.NewDummyLogger(), nil)

		httpReq := testutils.NewMockJSONRequest(t, test.Req)
		httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

		resp := h.postSendCoin(httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	}

}

func TestPostSendCoin_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name    string
		Err     error
		ExpCode int
	}{
		{"ToUser doesn't exist", domain.ErrUserNotFound, http.StatusBadRequest},
		{"SelfSending", domain.ErrBadRequest, http.StatusBadRequest},
		{"Low balance", domain.ErrLowBalance, http.StatusBadRequest},
		{"Unexpected DBError", errors.New("unexpected DBError"), http.StatusInternalServerError},
	}

	for _, test := range tests {
		svc := mocks.NewTransaction(t)
		h := NewTransactionHandler(testutils.NewDummyLogger(), svc)

		uID := 2
		req := types.PostSendCoinRequest{ToUser: "Avito", Amount: 100}
		httpReq := testutils.NewMockJSONRequest(t, req)
		httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

		svc.On("SendCoinByName",
			mock.Anything,
			domain.Transaction{From: uID, Amount: req.Amount},
			req.ToUser).
			Return(test.Err)

		resp := h.postSendCoin(httpReq)

		require.Equal(t, test.ExpCode, resp.StatusCode())
		svc.AssertExpectations(t)
	}
}

func TestGetBuyItemCoin_Success(t *testing.T) {
	t.Parallel()

	svc := mocks.NewTransaction(t)
	h := NewTransactionHandler(testutils.NewDummyLogger(), svc)

	uID := 2
	req := types.GetBuyItemRequest{
		Item: "AvitoHoody",
	}

	httpReq := testutils.NewMockRequestWithItemQueryVal(req.Item)
	httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

	svc.On(
		"BuyItemByName", mock.Anything, uID, req.Item).
		Return(nil)

	resp := h.getBuyItem(httpReq)

	require.Equal(t, http.StatusOK, resp.StatusCode())
	svc.AssertExpectations(t)
}

func TestGetBuyItemCoin_EmptyContextVal(t *testing.T) {
	t.Parallel()

	h := NewTransactionHandler(testutils.NewDummyLogger(), nil)

	req := types.GetBuyItemRequest{
		Item: "AvitoHoody",
	}

	httpReq := testutils.NewMockRequestWithItemQueryVal(req.Item)

	resp := h.getBuyItem(httpReq)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
}

func TestGetBuyItemCoin_BadRequestCases(t *testing.T) {
	t.Parallel()

	uID := 2
	tests := []struct {
		Name string
		Req  string
	}{
		{"Empty Query", ""},
	}

	for _, test := range tests {
		h := NewTransactionHandler(testutils.NewDummyLogger(), nil)

		httpReq := testutils.NewMockRequestWithItemQueryVal(test.Req)
		httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

		resp := h.getBuyItem(httpReq)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	}

}

func TestGetBuyItemCoin_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name    string
		Err     error
		ExpCode int
	}{
		{"Item doesn't exist or was deleted", domain.ErrMerchNotFound, http.StatusBadRequest},
		{"Low balance", domain.ErrLowBalance, http.StatusBadRequest},
		{"Unexpected DBError", errors.New("unexpected DBError"), http.StatusInternalServerError},
	}

	for _, test := range tests {
		svc := mocks.NewTransaction(t)
		h := NewTransactionHandler(testutils.NewDummyLogger(), svc)

		uID := 2
		req := types.GetBuyItemRequest{Item: "AvitoHoody"}
		httpReq := testutils.NewMockRequestWithItemQueryVal(req.Item)
		httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

		svc.On("BuyItemByName", mock.Anything, uID, req.Item).
			Return(test.Err)

		resp := h.getBuyItem(httpReq)

		require.Equal(t, test.ExpCode, resp.StatusCode())
		svc.AssertExpectations(t)
	}
}
