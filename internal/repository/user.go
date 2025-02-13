package repository

import (
	"avito_shop/internal/domain"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=User --filename=user_repository_mock.go
type User interface {
	Put(ctx context.Context, user domain.User) (domain.UserID, error)
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
	GetByName(ctx context.Context, name domain.UserName) (domain.User, error)
	GetInfoByID(ctx context.Context, id domain.UserID) (domain.UserInfo, error)
}
