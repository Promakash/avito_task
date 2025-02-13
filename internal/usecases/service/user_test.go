package service

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPut_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	user := domain.User{
		Name: "AvitoEmployee",
	}
	uid := 2
	ctx := context.Background()

	userRepo.On("Put", mock.Anything, user).
		Return(uid, nil)

	resultID, err := svc.Put(ctx, user)

	require.NoError(t, err)
	require.Equal(t, uid, resultID)
	userRepo.AssertExpectations(t)
}

func TestPut_AlreadyExists(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	user := domain.User{
		Name: "AvitoEmployee",
	}
	ctx := context.Background()

	userRepo.On("Put", mock.Anything, user).
		Return(0, domain.ErrUserExists)

	_, err := svc.Put(ctx, user)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUserExists)
	userRepo.AssertExpectations(t)
}

func TestPut_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	user := domain.User{
		Name: "AvitoEmployee",
	}
	ctx := context.Background()

	userRepo.On("Put", mock.Anything, user).
		Return(0, errors.New("new DBError"))

	_, err := svc.Put(ctx, user)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	user := domain.User{
		ID:   uid,
		Name: "AvitoEmployee",
		Info: domain.UserInfo{
			Coins: 100,
		},
	}
	ctx := context.Background()

	userRepo.On("GetByID", mock.Anything, uid).
		Return(user, nil)

	resultUser, err := svc.GetByID(ctx, user.ID)

	require.NoError(t, err)
	require.Equal(t, user, resultUser)
	userRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	user := domain.User{}
	ctx := context.Background()

	userRepo.On("GetByID", mock.Anything, uid).
		Return(user, domain.ErrUserNotFound)

	_, err := svc.GetByID(ctx, uid)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	userRepo.AssertExpectations(t)
}

func TestGetByID_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	user := domain.User{}
	ctx := context.Background()

	userRepo.On("GetByID", mock.Anything, uid).
		Return(user, errors.New("excellent error from DB"))

	_, err := svc.GetByID(ctx, uid)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestGetByName_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	name := "AvitoEmployee"
	user := domain.User{
		ID:   2,
		Name: name,
		Info: domain.UserInfo{
			Coins: 100,
		},
	}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, name).
		Return(user, nil)

	resultUser, err := svc.GetByName(ctx, name)

	require.NoError(t, err)
	require.Equal(t, user, resultUser)
	userRepo.AssertExpectations(t)
}

func TestGetByName_NotFound(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	name := "AvitoEmployee"
	user := domain.User{}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, name).
		Return(user, domain.ErrUserNotFound)

	_, err := svc.GetByName(ctx, name)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	userRepo.AssertExpectations(t)
}

func TestGetByName_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	name := "AvitoEmployee"
	user := domain.User{}
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, name).
		Return(user, errors.New("new avito db error"))

	_, err := svc.GetByName(ctx, name)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestGetInfoByID_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	uInfo := domain.UserInfo{Coins: 100}
	ctx := context.Background()

	userRepo.On("GetInfoByID", mock.Anything, uid).
		Return(uInfo, nil)

	resultInfo, err := svc.GetInfoByID(ctx, uid)

	require.NoError(t, err)
	require.Equal(t, uInfo, resultInfo)
	userRepo.AssertExpectations(t)
}

func TestGetInfoByID_NotFound(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	uInfo := domain.UserInfo{}
	ctx := context.Background()

	userRepo.On("GetInfoByID", mock.Anything, uid).
		Return(uInfo, domain.ErrUserNotFound)

	_, err := svc.GetInfoByID(ctx, uid)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	userRepo.AssertExpectations(t)
}

func TestGetInfoByID_DBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewUser(userRepo)

	uid := 2
	uInfo := domain.UserInfo{}
	ctx := context.Background()

	userRepo.On("GetInfoByID", mock.Anything, uid).
		Return(uInfo, errors.New("error from db"))

	_, err := svc.GetInfoByID(ctx, uid)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}
