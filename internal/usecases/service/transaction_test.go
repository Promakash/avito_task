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

func TestSendCoinByName_Success(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	TxRepo := mocks.NewTransaction(t)
	svc := NewTransaction(TxRepo, userRepo, nil)

	fromID := 0
	toID := 1
	amount := 100
	toName := "Avito"
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, toName).
		Return(domain.User{ID: toID, Name: toName}, nil)
	TxRepo.On(
		"SendCoin",
		mock.Anything,
		domain.Transaction{From: fromID, To: toID, Amount: amount}).
		Return(nil)

	err := svc.SendCoinByName(ctx, domain.Transaction{From: fromID, Amount: amount}, toName)

	require.NoError(t, err)
	userRepo.AssertExpectations(t)
	TxRepo.AssertExpectations(t)
}

func TestSendCoinByName_SelfSending(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewTransaction(nil, userRepo, nil)

	fromID := 0
	toID := fromID
	amount := 100
	toName := "Avito"
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, toName).
		Return(domain.User{ID: toID, Name: toName}, nil)

	err := svc.SendCoinByName(ctx, domain.Transaction{From: fromID, Amount: amount}, toName)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrBadRequest)
	userRepo.AssertExpectations(t)
}

func TestSendCoinByName_InvalidUsername(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewTransaction(nil, userRepo, nil)

	fromID := 0
	amount := 100
	toName := "Avito"
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, toName).
		Return(domain.User{}, domain.ErrUserNotFound)

	err := svc.SendCoinByName(ctx, domain.Transaction{From: fromID, Amount: amount}, toName)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	userRepo.AssertExpectations(t)
}

func TestSendCoinByName_CheckingUsernameDBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	svc := NewTransaction(nil, userRepo, nil)

	fromID := 0
	amount := 100
	toName := "Avito"
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, toName).
		Return(domain.User{}, errors.New("unexpected DBError"))

	err := svc.SendCoinByName(ctx, domain.Transaction{From: fromID, Amount: amount}, toName)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSendCoinByName_MakingTxDBError(t *testing.T) {
	t.Parallel()

	userRepo := mocks.NewUser(t)
	TxRepo := mocks.NewTransaction(t)
	svc := NewTransaction(TxRepo, userRepo, nil)

	fromID := 0
	toID := 1
	amount := 100
	toName := "Avito"
	ctx := context.Background()

	userRepo.On("GetByName", mock.Anything, toName).
		Return(domain.User{ID: toID, Name: toName}, nil)
	TxRepo.On(
		"SendCoin",
		mock.Anything,
		domain.Transaction{From: fromID, To: toID, Amount: amount}).
		Return(errors.New("unexpected DB error"))

	err := svc.SendCoinByName(ctx, domain.Transaction{From: fromID, Amount: amount}, toName)

	require.Error(t, err)
	userRepo.AssertExpectations(t)
	TxRepo.AssertExpectations(t)
}

func TestBuyItemByName_Success(t *testing.T) {
	t.Parallel()

	TxRepo := mocks.NewTransaction(t)
	merchRepo := mocks.NewMerch(t)
	svc := NewTransaction(TxRepo, nil, merchRepo)

	buyerID := 0
	merch := domain.Merch{
		ID:    1,
		Name:  "AvitoHoody",
		Price: 100,
	}
	ctx := context.Background()

	merchRepo.On("GetByName", mock.Anything, merch.Name).
		Return(merch, nil)
	TxRepo.On("BuyItem", mock.Anything, buyerID, merch).
		Return(nil)

	err := svc.BuyItemByName(ctx, buyerID, merch.Name)

	require.NoError(t, err)
	TxRepo.AssertExpectations(t)
}

func TestBuyItemByName_InvalidMerchName(t *testing.T) {
	t.Parallel()

	merchRepo := mocks.NewMerch(t)
	svc := NewTransaction(nil, nil, merchRepo)

	buyerID := 0
	merch := domain.Merch{
		Name: "InvalidHoody",
	}
	ctx := context.Background()

	merchRepo.On("GetByName", mock.Anything, merch.Name).
		Return(domain.Merch{}, domain.ErrMerchNotFound)

	err := svc.BuyItemByName(ctx, buyerID, merch.Name)

	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrMerchNotFound)
}

func TestBuyItemByName_CheckingMerchNameDBError(t *testing.T) {
	t.Parallel()

	merchRepo := mocks.NewMerch(t)
	svc := NewTransaction(nil, nil, merchRepo)

	buyerID := 0
	merch := domain.Merch{
		ID:    1,
		Name:  "AvitoHoody",
		Price: 100,
	}
	ctx := context.Background()

	merchRepo.On("GetByName", mock.Anything, merch.Name).
		Return(domain.Merch{}, errors.New("excellent DBError"))

	err := svc.BuyItemByName(ctx, buyerID, merch.Name)

	require.Error(t, err)
}

func TestBuyItemByName_PurchaseDBError(t *testing.T) {
	t.Parallel()

	TxRepo := mocks.NewTransaction(t)
	merchRepo := mocks.NewMerch(t)
	svc := NewTransaction(TxRepo, nil, merchRepo)

	buyerID := 0
	merch := domain.Merch{
		ID:    1,
		Name:  "AvitoHoody",
		Price: 100,
	}
	ctx := context.Background()

	merchRepo.On("GetByName", mock.Anything, merch.Name).
		Return(merch, nil)
	TxRepo.On("BuyItem", mock.Anything, buyerID, merch).
		Return(errors.New("cool error - tx is down"))

	err := svc.BuyItemByName(ctx, buyerID, merch.Name)

	require.Error(t, err)
	TxRepo.AssertExpectations(t)
}
