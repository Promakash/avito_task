package tests

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/pkg/testutils"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const buyItemPath = "/buy"
const buyItemMethod = http.MethodGet

func buyItemHelper(t *testing.T, req types.GetBuyItemRequest, token string) *http.Response {
	path := fmt.Sprintf("%s%s/%s", apiPath, buyItemPath, req.Item)
	resp, err := testutils.SendRequest(t, path, buyItemMethod, token, nil)

	require.NoError(t, err)

	return resp
}

func TestGetBuyItem_Success(t *testing.T) {
	userCreds := types.PostAuthRequest{
		Username: "AvitoBuyItem",
		Password: "12345",
	}
	token := getTokenHelper(t, userCreds)

	req := types.GetBuyItemRequest{Item: "pen"}
	expStatus := http.StatusOK

	resp := buyItemHelper(t, req, token)

	require.Equal(t, expStatus, resp.StatusCode)
}

func TestGetBuyItem_BadRequest(t *testing.T) {
	userCreds := types.PostAuthRequest{
		Username: "AvitoBuyItem",
		Password: "12345",
	}
	token := getTokenHelper(t, userCreds)

	expStatus := http.StatusBadRequest
	req := types.GetBuyItemRequest{Item: "spaceship"}

	resp := buyItemHelper(t, req, token)

	require.Equal(t, expStatus, resp.StatusCode)
}

func TestGetBuyItem_UnauthorizedCases(t *testing.T) {
	req := types.GetBuyItemRequest{
		Item: "powerbank",
	}
	path := fmt.Sprintf("%s%s/%s", apiPath, buyItemPath, req.Item)
	expStatus := http.StatusUnauthorized

	tests := []struct {
		Name  string
		Token string
	}{
		{"Empty token", ""},
		{"Invalid token", "32135dsvcxa"},
	}

	for _, test := range tests {
		resp, err := testutils.SendRequest(t, path, buyItemMethod, test.Token, nil)
		require.NoError(t, err)

		require.Equal(t, expStatus, resp.StatusCode)
	}
}

func TestGetBuyItem_MoneyEnded(t *testing.T) {
	// By the task, default amount of coins of new user is 1000
	// So we will buy 3 pink-hoody (price 500). Two Ok's and 1 bad request expected

	userCreds := types.PostAuthRequest{Username: "AvitoBrokeBuyer", Password: "12345"}
	token := getTokenHelper(t, userCreds)

	req := types.GetBuyItemRequest{Item: "pink-hoody"}
	expStatuses := []int{http.StatusOK, http.StatusOK, http.StatusBadRequest}

	for _, exp := range expStatuses {
		resp := buyItemHelper(t, req, token)

		require.Equal(t, exp, resp.StatusCode)
	}
}
