package service

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"avito_shop/internal/usecases"
	"context"
	"fmt"
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
	uid, err := s.repo.Put(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("UserService.Put: %w", err)
	}

	return uid, nil
}

func (s *User) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("UserService.GetByID: %w", err)
	}

	return user, nil
}

func (s *User) GetByName(ctx context.Context, name domain.UserName) (domain.User, error) {
	user, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return domain.User{}, fmt.Errorf("UserService.GetByName: %w", err)
	}

	return user, nil
}

func (s *User) GetInfoByID(ctx context.Context, id domain.UserID) (domain.UserInfo, error) {
	info, err := s.repo.GetInfoByID(ctx, id)
	if err != nil {
		return domain.UserInfo{}, fmt.Errorf("UserService.GetInfoByID: %w", err)
	}

	return info, nil
}
