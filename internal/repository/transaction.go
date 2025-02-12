package repository

import (
	"avito_shop/internal/domain"
	"context"
)

const ShopDBID = 1

type Transaction interface {
	Put(ctx context.Context, tx domain.Transaction) error
	PurchaseItem(ctx context.Context, tx domain.Transaction) error
}
