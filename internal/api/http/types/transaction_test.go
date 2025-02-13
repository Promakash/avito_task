package types_test

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/lib/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePostSendCoinRequest_Success(t *testing.T) {
	t.Parallel()

	req := &types.PostSendCoinRequest{
		ToUser: "Avito",
		Amount: 100,
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	result, err := types.CreatePostSendCoinRequest(httpReq)

	require.NoError(t, err)
	require.Equal(t, req, result)
}

func TestCreatePostSendCoinRequest_EmptyToUser(t *testing.T) {
	t.Parallel()

	req := &types.PostSendCoinRequest{
		Amount: 100,
	}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostSendCoinRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostSendCoinRequest_AmountZeroValue(t *testing.T) {
	t.Parallel()

	req := &types.PostSendCoinRequest{ToUser: "Avito"}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostSendCoinRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostSendCoinRequest_EmptyReq(t *testing.T) {
	t.Parallel()

	req := &types.PostSendCoinRequest{}

	httpReq := testutils.NewMockJSONRequest(t, req)

	_, err := types.CreatePostSendCoinRequest(httpReq)

	require.Error(t, err)
}

func TestCreatePostSendCoinRequest_BrokenJSON(t *testing.T) {
	t.Parallel()

	brokenJSON := []byte("{\"toUser\":\"avito\",\"amount\":\"100\"")

	httpReq := testutils.NewMockJSONRequest(t, brokenJSON)

	_, err := types.CreatePostSendCoinRequest(httpReq)

	require.Error(t, err)
}

func TestCreateGetBuyItemRequest_Success(t *testing.T) {
	itemName := "powerbank"

	httpReq := testutils.NewMockRequestWithItemQueryVal(itemName)

	result, err := types.CreateGetBuyItemRequest(httpReq)

	require.NoError(t, err)
	require.Equal(t, itemName, result.Item)
}

func TestCreateGetBuyItemRequest_QueryZeroVal(t *testing.T) {
	itemName := ""

	httpReq := testutils.NewMockRequestWithItemQueryVal(itemName)

	_, err := types.CreateGetBuyItemRequest(httpReq)

	require.Error(t, err)
}
