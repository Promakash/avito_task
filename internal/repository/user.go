package repository

import (
	"avito_shop/internal/domain"
	"context"
)

type User interface {
	Put(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
	GetByName(ctx context.Context, name domain.UserName) (domain.User, error)
	GetInfoByID(ctx context.Context, id domain.UserID) (domain.UserInfo, error)
}
