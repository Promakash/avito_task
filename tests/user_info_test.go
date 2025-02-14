package tests

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	"avito_shop/internal/usecases/service"
	"avito_shop/pkg/testutils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const userInfoPath = "/info"
const userInfoMethod = http.MethodGet

func userInfoHelper(t *testing.T, token string) *http.Response {
	path := fmt.Sprintf("%s%s", apiPath, userInfoPath)

	resp, err := testutils.SendRequest(t, path, userInfoMethod, token, nil)
	require.NoError(t, err)

	return resp
}

func TestGetInfo_Success(t *testing.T) {
	userCreds := types.PostAuthRequest{
		Username: "AvitoHappyEmployee",
		Password: "12345",
	}

	token := getTokenHelper(t, userCreds)
	expStatus := http.StatusOK
	expResult := types.GetInfoResponse{Coins: 1000}

	resp := userInfoHelper(t, token)
	require.Equal(t, expStatus, resp.StatusCode)

	var payload types.GetInfoResponse
	err := json.NewDecoder(resp.Body).Decode(&payload)
	require.NoError(t, err)

	require.Equal(t, expResult, payload)
}

func TestGetInfo_DeletedUserRequest(t *testing.T) {
	authSecret := os.Getenv("AUTH_SECRET")
	require.NotEmpty(t, authSecret)

	delUser := domain.User{
		ID:   2415,
		Name: "DELETED",
	}

	authService := service.NewAuth(nil, authSecret)
	authToken, err := authService.GenerateToken(delUser)
	require.NoError(t, err)

	expStatus := http.StatusBadRequest

	resp := userInfoHelper(t, authToken)

	require.Equal(t, expStatus, resp.StatusCode)
}

func TestGetInfo_UnauthorizedCases(t *testing.T) {
	expStatus := http.StatusUnauthorized

	tests := []struct {
		Name  string
		Token string
	}{
		{"Empty token", ""},
		{"Invalid token", "32135dsvcxa"},
	}

	for _, test := range tests {
		resp := userInfoHelper(t, test.Token)

		require.Equal(t, expStatus, resp.StatusCode)
	}
}
