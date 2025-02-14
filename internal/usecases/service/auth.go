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
		return "", fmt.Errorf("AuthService.Login: %w", err)
	}

	if ok := s.compareHash(user.HashedPassword, password); ok {
		return s.GenerateToken(user)
	}

	return "", fmt.Errorf("AuthService.Login: different password hash: %w", domain.ErrUnauthorized)
}

func (s *Auth) Register(ctx context.Context, username domain.UserName, password string) (domain.Token, error) {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return "", fmt.Errorf("AuthService.Register: %w", err)
	}

	user := domain.User{
		Name:           username,
		HashedPassword: hashedPassword,
	}

	id, err := s.userRepo.Put(ctx, user)
	if err != nil {
		return "", fmt.Errorf("AuthService.Register: %w", err)
	}
	user.ID = id

	return s.GenerateToken(user)
}

func (s *Auth) GenerateToken(user domain.User) (domain.Token, error) {
	if user.ID == 0 {
		return "", fmt.Errorf("AuthService.GenerateToken: %w", errors.New("invalid id"))
	}

	token, err := libjwt.NewToken(user, s.secret)
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateToken: %w", err)
	}

	return token, nil
}

func (s *Auth) ParseToken(token domain.Token) (domain.UserID, error) {
	var id domain.UserID

	val, err := libjwt.ParseToken(token, s.secret)
	if err != nil {
		return id, fmt.Errorf("AuthService.ParseToken: %w", err)
	}

	return val, nil
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
	return err == nil
}
