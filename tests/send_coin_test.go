package tests

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/pkg/testutils"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const sendCoinPath = "/sendCoin"
const sendCoinMethod = http.MethodPost

func sendCoinHelper(t *testing.T, req types.PostSendCoinRequest, token string) *http.Response {
	path := fmt.Sprintf("%s%s", apiPath, sendCoinPath)

	resp, err := testutils.SendRequest(t, path, sendCoinMethod, token, &req)
	require.NoError(t, err)

	return resp
}

func TestPostSendCoin_Success(t *testing.T) {
	sender := types.PostAuthRequest{
		Username: "AvitoSender",
		Password: "12345",
	}
	receiver := types.PostAuthRequest{
		Username: "AvitoReceiver",
		Password: "12345",
	}

	token := getTokenHelper(t, sender)
	_ = getTokenHelper(t, receiver) // need to create receiver

	req := types.PostSendCoinRequest{ToUser: receiver.Username, Amount: 1}
	expStatus := http.StatusOK

	resp := sendCoinHelper(t, req, token)

	require.Equal(t, expStatus, resp.StatusCode)
}

func TestPostSendCoin_BadRequestCases(t *testing.T) {
	sender := types.PostAuthRequest{
		Username: "AvitoSender",
		Password: "12345",
	}
	receiver := types.PostAuthRequest{
		Username: "AvitoReceiver",
		Password: "12345",
	}

	token := getTokenHelper(t, sender)
	_ = getTokenHelper(t, receiver) // need to create receiver
	expStatus := http.StatusBadRequest

	tests := []struct {
		Name string
		Req  types.PostSendCoinRequest
	}{
		{"Empty toUser", types.PostSendCoinRequest{Amount: 100}},
		{"Empty amount", types.PostSendCoinRequest{ToUser: receiver.Username}},
		{"Empty Request", types.PostSendCoinRequest{}},
	}

	for _, test := range tests {
		resp := sendCoinHelper(t, test.Req, token)

		require.Equal(t, expStatus, resp.StatusCode)
	}

	path := fmt.Sprintf("%s%s", apiPath, sendCoinPath)
	brokenJSON := []byte("{\"toUser\":\"AvitoReceiver\",\"amount\":\"100\"")

	resp, err := testutils.SendRequest(t, path, sendCoinMethod, token, &brokenJSON)
	require.NoError(t, err)

	require.Equal(t, expStatus, resp.StatusCode)
}

func TestPostSendCoin_UnauthorizedCases(t *testing.T) {
	req := types.PostSendCoinRequest{
		ToUser: "AvitoReceiver",
		Amount: 1,
	}
	expStatus := http.StatusUnauthorized

	tests := []struct {
		Name  string
		Token string
	}{
		{"Empty token", ""},
		{"Invalid token", "32135dsvcxa"},
	}

	for _, test := range tests {
		resp := sendCoinHelper(t, req, test.Token)

		require.Equal(t, expStatus, resp.StatusCode)
	}
}

func TestPostSendCoin_MoneyEnded(t *testing.T) {
	// By the task, default amount of coins of new user is 1000
	// So we will try to send 1000 coins two times. First OK and second bad request expected

	sender := types.PostAuthRequest{
		Username: "AvitoBrokeSender",
		Password: "12345",
	}
	receiver := types.PostAuthRequest{
		Username: "AvitoRichReceiver",
		Password: "12345",
	}

	token := getTokenHelper(t, sender)
	_ = getTokenHelper(t, receiver) // need to create receiver

	req := types.PostSendCoinRequest{ToUser: receiver.Username, Amount: 1000}
	expStatuses := []int{http.StatusOK, http.StatusBadRequest}

	for _, exp := range expStatuses {
		resp := sendCoinHelper(t, req, token)

		require.Equal(t, exp, resp.StatusCode)
	}
}
