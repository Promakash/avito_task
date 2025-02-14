package service

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const secretForTests string = "secret"

func hashPasswordForTests(t *testing.T, password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return hashedPassword
}

func TestLogin_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewAuth(userRepo, secretForTests)

	password := "12345"
	user := domain.User{
		ID:             2,
		Name:           "Avito",
		HashedPassword: hashPasswordForTests(t, password),
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, user.Name).
		Return(user, nil)

	_, err := svc.Login(ctx, user.Name, password)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewAuth(userRepo, secretForTests)

	password := "12345"
	user := domain.User{
		ID:             2,
		Name:           "Avito",
		HashedPassword: []byte("Nope bad hash"),
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, user.Name).
		Return(user, nil)

	_, err := svc.Login(ctx, user.Name, password)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUnauthorized)
	userRepo.AssertExpectations(t)
}

func TestLogin_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewAuth(userRepo, secretForTests)

	password := "12345"
	user := domain.User{
		Name: "Avito",
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, user.Name).
		Return(domain.User{}, errors.New("db err"))

	_, err := svc.Login(ctx, user.Name, password)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewAuth(userRepo, secretForTests)

	password := "12345"
	uid := 2
	user := domain.User{
		Name: "Avito",
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, user.Name).
		Return(domain.User{}, domain.ErrUserNotFound)
	userRepo.On("Put", mock.Anything, mock.Anything).
		Return(uid, nil)

	_, err := svc.Login(ctx, user.Name, password)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestRegister_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewAuth(userRepo, secretForTests)

	password := "12345"
	uid := 0
	user := domain.User{
		Name: "Avito",
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, user.Name).
		Return(domain.User{}, domain.ErrUserNotFound)
	userRepo.On("Put", mock.Anything, mock.Anything).
		Return(uid, errors.New("db err"))

	_, err := svc.Login(ctx, user.Name, password)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestGenerateToken_Success(t *testing.T) {
	t.Parallel()

	svc := NewAuth(nil, secretForTests)

	user := domain.User{
		ID:   2,
		Name: "Avito",
	}

	_, err := svc.GenerateToken(user)

	require.NoError(t, err)
}

func TestGenerateToken_IncorrectID(t *testing.T) {
	t.Parallel()

	svc := NewAuth(nil, secretForTests)

	user := domain.User{}

	_, err := svc.GenerateToken(user)

	require.Error(t, err)
}

func TestParseToken_Success(t *testing.T) {
	t.Parallel()

	svc := NewAuth(nil, secretForTests)

	user := domain.User{
		ID:   2,
		Name: "Avito",
	}

	token, err := svc.GenerateToken(user)
	require.NoError(t, err)

	parsedVal, err := svc.ParseToken(token)

	require.NoError(t, err)
	require.Equal(t, user.ID, parsedVal)
}

func TestParseToken_IncorrectToken(t *testing.T) {
	t.Parallel()

	svc := NewAuth(nil, secretForTests)

	token := "ddasdxbe1x9g5z"

	_, err := svc.ParseToken(token)

	require.Error(t, err)
}
