package usecases

import (
	"avito_shop/internal/domain"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=Transaction --filename=tx_service_mock.go
type Transaction interface {
	SendCoinByName(ctx context.Context, tx domain.Transaction, to domain.UserName) error
	BuyItemByName(ctx context.Context, uid domain.UserID, name domain.MerchName) error
}
