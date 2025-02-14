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

func TestGetInfo_Success(t *testing.T) {
	t.Parallel()

	svc := mocks.NewUser(t)
	h := NewUserHandler(testutils.NewDummyLogger(), svc)

	uID := 2
	httpReq := testutils.NewMockRequest()
	httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)
	expected := &types.GetInfoResponse{Coins: 100}

	svc.On("GetInfoByID", mock.Anything, uID).
		Return(domain.UserInfo{Coins: expected.Coins}, nil)

	resp := h.getInfo(httpReq)

	require.Equal(t, http.StatusOK, resp.StatusCode())
	require.Equal(t, expected, resp.GetPayload())
	svc.AssertExpectations(t)
}

func TestGetInfo_EmptyContextVal(t *testing.T) {
	t.Parallel()

	h := NewUserHandler(testutils.NewDummyLogger(), nil)

	httpReq := testutils.NewMockRequest()

	resp := h.getInfo(httpReq)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
}

func TestGetInfo_UserDeletedBefore(t *testing.T) {
	t.Parallel()

	svc := mocks.NewUser(t)
	h := NewUserHandler(testutils.NewDummyLogger(), svc)

	uID := 2
	httpReq := testutils.NewMockRequest()
	httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

	svc.On("GetInfoByID", mock.Anything, uID).
		Return(domain.UserInfo{}, domain.ErrUserNotFound)

	resp := h.getInfo(httpReq)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode())
	svc.AssertExpectations(t)
}

func TestGetInfo_TxFailed(t *testing.T) {
	t.Parallel()

	svc := mocks.NewUser(t)
	h := NewUserHandler(testutils.NewDummyLogger(), svc)

	uID := 2
	httpReq := testutils.NewMockRequest()
	httpReq = testutils.AddUserIDToRequestContext(httpReq, uID)

	svc.On("GetInfoByID", mock.Anything, uID).
		Return(domain.UserInfo{}, errors.New("tx failed - rolled back"))

	resp := h.getInfo(httpReq)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	svc.AssertExpectations(t)
}
