package usecases

import (
	"avito_shop/internal/domain"
	"context"
)

type Auth interface {
	Login(ctx context.Context, username domain.UserName, password string) (domain.Token, error)
	Register(ctx context.Context, username domain.UserName, password string) (domain.Token, error)
	GenerateToken(user domain.User) (domain.Token, error)
	ParseToken(token domain.Token) (domain.UserID, error)
}
