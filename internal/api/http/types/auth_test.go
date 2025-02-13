package types_test

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/lib/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePostAuthRequest_Success(t *testing.T) {
	t.Parallel()

	req := &types.PostAuthRequest{
		Username: "Avito",
		Password: "12345",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	result, err := types.CreatePostAuthRequest(httpReq)

	require.NoError(t, err)
	require.Equal(t, req, result)
}

func TestCreatePostAuthRequest_EmptyUsername(t *testing.T) {
	t.Parallel()

	req := &types.PostAuthRequest{
		Password: "12345",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_EmptyPassword(t *testing.T) {
	t.Parallel()

	req := &types.PostAuthRequest{
		Username: "Avito",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_EmptyReq(t *testing.T) {
	t.Parallel()

	req := &types.PostAuthRequest{}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_BrokenJSON(t *testing.T) {
	t.Parallel()

	brokenJSON := []byte("{\"username\":\"avito\",\"password\":\"12345\"")

	httpReq := testutils.NewMockJSONRequest(t, brokenJSON)

	_, err := types.CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}
