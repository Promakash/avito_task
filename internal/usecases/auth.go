package usecases

import (
	"avito_shop/internal/domain"
	"context"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=Auth --filename=auth_service_mock.go
type Auth interface {
	Login(ctx context.Context, username domain.UserName, password string) (domain.Token, error)
	Register(ctx context.Context, username domain.UserName, password string) (domain.Token, error)
	GenerateToken(user domain.User) (domain.Token, error)
	ParseToken(token domain.Token) (domain.UserID, error)
}
