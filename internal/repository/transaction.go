package repository

import (
	"avito_shop/internal/domain"
	"context"
)

const ShopDBID = 1

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=Transaction --filename=tx_repository_mock.go
type Transaction interface {
	SendCoin(ctx context.Context, tx domain.Transaction) error
	BuyItem(ctx context.Context, uid domain.UserID, item domain.Merch) error
}
