package types

import (
	"avito_shop/pkg/testutils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePostAuthRequest_Success(t *testing.T) {
	t.Parallel()

	req := &PostAuthRequest{
		Username: "Avito",
		Password: "12345",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	result, err := CreatePostAuthRequest(httpReq)

	require.NoError(t, err)
	require.Equal(t, req, result)
}

func TestCreatePostAuthRequest_EmptyUsername(t *testing.T) {
	t.Parallel()

	req := &PostAuthRequest{
		Password: "12345",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_EmptyPassword(t *testing.T) {
	t.Parallel()

	req := &PostAuthRequest{
		Username: "Avito",
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_EmptyReq(t *testing.T) {
	t.Parallel()

	req := &PostAuthRequest{}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostAuthRequest_BrokenJSON(t *testing.T) {
	t.Parallel()

	brokenJSON := []byte("{\"username\":\"avito\",\"password\":\"12345\"")

	httpReq := testutils.NewMockJSONRequest(t, brokenJSON)

	_, err := CreatePostAuthRequest(httpReq)

	require.Error(t, err)
}
