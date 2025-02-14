package types

import "avito_shop/internal/domain"

type GetInfoResponse struct {
	Coins       int                `json:"coins"`
	Inventory   []domain.Inventory `json:"inventory"`
	CoinHistory CoinHistory        `json:"coinHistory"`
}

type CoinHistory struct {
	Received []CoinHistoryIncoming `json:"received"`
	Sent     []CoinHistoryUpcoming `json:"sent"`
}

type CoinHistoryIncoming struct {
	FromUser domain.UserName `json:"fromUser"`
	Amount   int             `json:"amount"`
}

type CoinHistoryUpcoming struct {
	ToUser domain.UserName `json:"ToUser"`
	Amount int             `json:"amount"`
}

func CreateGetInfoResponse(info domain.UserInfo) *GetInfoResponse {
	var incoming []CoinHistoryIncoming
	var upcoming []CoinHistoryUpcoming

	for _, tx := range info.Transactions {
		if tx.Direction == domain.Received {
			recTx := CoinHistoryIncoming{
				FromUser: tx.OtherUser,
				Amount:   tx.Amount,
			}
			incoming = append(incoming, recTx)
		} else {
			sentTx := CoinHistoryUpcoming{
				ToUser: tx.OtherUser,
				Amount: tx.Amount,
			}
			upcoming = append(upcoming, sentTx)
		}
	}

	return &GetInfoResponse{
		Coins:     info.Coins,
		Inventory: info.Inventory,
		CoinHistory: CoinHistory{
			Received: incoming,
			Sent:     upcoming,
		},
	}
}
