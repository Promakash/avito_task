package service

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"avito_shop/internal/usecases"
	"context"
)

type User struct {
	repo repository.User
}

func NewUser(repo repository.User) usecases.User {
	return &User{
		repo: repo,
	}
}

func (s *User) Put(ctx context.Context, user domain.User) (domain.UserID, error) {
	return s.repo.Put(ctx, user)
}

func (s *User) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *User) GetByName(ctx context.Context, name domain.UserName) (domain.User, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *User) GetInfoByID(ctx context.Context, id domain.UserID) (domain.UserInfo, error) {
	return s.repo.GetInfoByID(ctx, id)
}
