package service

import (
	"avito_shop/internal/domain"
	libjwt "avito_shop/internal/lib/jwt"
	"avito_shop/internal/repository"
	"avito_shop/internal/usecases"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"context"
)

type Auth struct {
	userRepo repository.User
	secret   string
}

func NewAuth(repo repository.User, secret string) usecases.Auth {
	return &Auth{
		userRepo: repo,
		secret:   secret,
	}
}

func (s *Auth) Login(ctx context.Context, username domain.UserName, password string) (domain.Token, error) {
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return s.Register(ctx, username, password)
		}
		return "", err
	}

	if ok := s.compareHash(user.HashedPassword, password); ok {
		return s.GenerateToken(user)
	}

	return "", fmt.Errorf("different password hash: %w", domain.ErrUnauthorized)
}

func (s *Auth) Register(ctx context.Context, username domain.UserName, password string) (domain.Token, error) {
	hashedPasword, err := s.hashPassword(password)
	if err != nil {
		return "", err
	}

	user := domain.User{
		Name:           username,
		HashedPassword: hashedPasword,
	}

	err = s.userRepo.Put(ctx, user)
	if err != nil {
		return "", err
	}

	return s.GenerateToken(user)
}

func (s *Auth) GenerateToken(user domain.User) (domain.Token, error) {
	return libjwt.NewToken(user, s.secret)
}

func (s *Auth) ParseToken(token domain.Token) (domain.UserID, error) {
	var id domain.UserID

	val, err := libjwt.ParseToken(token, s.secret)
	if err != nil {
		return id, err
	}

	return val.(domain.UserID), nil
}

func (s *Auth) hashPassword(password string) (domain.UserHashPass, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (s *Auth) compareHash(hashedPassword domain.UserHashPass, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, domain.UserHashPass(password))
	if err != nil {
		return false
	}

	return true
}
