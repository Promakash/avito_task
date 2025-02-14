package tests

import (
	"avito_shop/internal/api/http/types"
	"avito_shop/internal/domain"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullUserFlow(t *testing.T) {
	// Testing full user's flow: Register -> buyItem -> sendCoin -> getInfoSender -> getInfoReceiver

	userCreds := types.PostAuthRequest{
		Username: "e2ebasicflow",
		Password: "12345",
	}

	token := getTokenHelper(t, userCreds)

	reqItem := types.GetBuyItemRequest{Item: "hoody"}

	// 1000 - 300 = 700
	itemResp := buyItemHelper(t, reqItem, token)
	require.Equal(t, http.StatusOK, itemResp.StatusCode)

	// Make new user for transfering
	receiver := types.PostAuthRequest{
		Username: "e2ebasicflowreceiver",
		Password: "12345",
	}
	recToken := getTokenHelper(t, receiver)

	reqTx := types.PostSendCoinRequest{
		ToUser: receiver.Username,
		Amount: 200,
	}

	// 700 - 200 = 500
	sendResp := sendCoinHelper(t, reqTx, token)
	require.Equal(t, http.StatusOK, sendResp.StatusCode)

	expSenderInfo := types.GetInfoResponse{
		Coins: 500,
		Inventory: []domain.Inventory{
			{Name: reqItem.Item, Quantity: 1},
		},
		CoinHistory: types.CoinHistory{
			Received: nil,
			Sent: []types.CoinHistoryUpcoming{
				{ToUser: receiver.Username, Amount: reqTx.Amount},
				{ToUser: "shop", Amount: 300},
			},
		},
	}

	infoResp := userInfoHelper(t, token)
	require.Equal(t, http.StatusOK, infoResp.StatusCode)

	var infoSenderPayload types.GetInfoResponse
	err := json.NewDecoder(infoResp.Body).Decode(&infoSenderPayload)
	require.NoError(t, err)

	require.Equal(t, expSenderInfo, infoSenderPayload)

	// Now checking receiver info
	expReceiverInfo := types.GetInfoResponse{
		Coins:     1200,
		Inventory: nil,
		CoinHistory: types.CoinHistory{
			Received: []types.CoinHistoryIncoming{
				{FromUser: userCreds.Username, Amount: reqTx.Amount},
			},
			Sent: nil,
		},
	}

	recInfoResp := userInfoHelper(t, recToken)
	require.Equal(t, http.StatusOK, recInfoResp.StatusCode)

	var infoRecPayload types.GetInfoResponse
	err = json.NewDecoder(recInfoResp.Body).Decode(&infoRecPayload)
	require.NoError(t, err)

	require.Equal(t, expReceiverInfo, infoRecPayload)
}
