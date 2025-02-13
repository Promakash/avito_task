package repository

import (
	"avito_shop/internal/domain"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=Merch --filename=merch_repository_mock.go
type Merch interface {
	GetByName(ctx context.Context, name string) (domain.Merch, error)
}
