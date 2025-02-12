package repository

import (
	"avito_shop/internal/domain"
	"context"
)

type Merch interface {
	GetByName(ctx context.Context, name string) (domain.Merch, error)
}
