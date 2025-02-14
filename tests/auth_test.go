package tests

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	pkgresp "avito_shop/pkg/http/responses"
	"avito_shop/pkg/testutils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var apiPath = os.Getenv("API_PATH")

const authPath = "/auth"
const postAuthMethod = http.MethodPost

func getTokenHelper(t *testing.T, userCreds types.PostAuthRequest) string {
	tokenPath := fmt.Sprintf("%s%s", apiPath, authPath)
	expStatus := http.StatusOK

	resp, err := testutils.SendRequest(t, tokenPath, postAuthMethod, "", &userCreds)
	require.NoError(t, err)
	require.Equal(t, expStatus, resp.StatusCode)

	var payload types.PostAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	require.NoError(t, err)

	require.NotEmpty(t, payload.Token)

	return payload.Token
}

func TestPostAuthHandler_UserDoesntExist(t *testing.T) {
	req := types.PostAuthRequest{
		Username: "AvitoAuth",
		Password: "12345",
	}

	_ = getTokenHelper(t, req)
}

func TestPostAuthHandler_LoginTwoTimes(t *testing.T) {
	req := types.PostAuthRequest{
		Username: "AvitoAuth",
		Password: "12345",
	}

	token1 := getTokenHelper(t, req)
	token2 := getTokenHelper(t, req)

	require.Equal(t, token1, token2)
}

func TestPostAuthHandler_InvalidPassword(t *testing.T) {
	path := fmt.Sprintf("%s%s", apiPath, authPath)

	req := types.PostAuthRequest{
		Username: "AvitoAuth",
		Password: "invalid password",
	}
	expStatus := http.StatusUnauthorized

	resp, err := testutils.SendRequest(t, path, postAuthMethod, "", &req)
	require.NoError(t, err)
	require.Equal(t, expStatus, resp.StatusCode)

	var errPayload pkgresp.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errPayload)
	require.NoError(t, err)

	require.Equal(t, errPayload.Message, domain.ErrUnauthorized.Error())
}

func TestPostAuthHandler_BadRequestCases(t *testing.T) {
	path := fmt.Sprintf("%s%s", apiPath, authPath)
	expStatus := http.StatusBadRequest

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
		resp, err := testutils.SendRequest(t, path, postAuthMethod, "", &test.Req)

		require.NoError(t, err)
		require.Equal(t, expStatus, resp.StatusCode)
	}
}
