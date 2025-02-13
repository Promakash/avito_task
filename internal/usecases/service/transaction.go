package service

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"avito_shop/internal/usecases"
	"context"
	"fmt"
)

type Transaction struct {
	repo      repository.Transaction
	userRepo  repository.User
	merchRepo repository.Merch
}

func NewTransaction(
	repo repository.Transaction,
	userRepo repository.User,
	merchRepo repository.Merch,
) usecases.Transaction {
	return &Transaction{
		repo:      repo,
		userRepo:  userRepo,
		merchRepo: merchRepo,
	}
}

func (s *Transaction) SendCoinByName(ctx context.Context, tx domain.Transaction, to domain.UserName) error {
	toUser, err := s.userRepo.GetByName(ctx, to)
	if err != nil {
		return fmt.Errorf("TxService.SendCoinByName: %w", err)
	}

	if toUser.ID == tx.From {
		return fmt.Errorf("TxService.SendCoinByName: can't allow selfsending: %w", domain.ErrBadRequest)
	}
	tx.To = toUser.ID

	err = s.repo.SendCoin(ctx, tx)
	if err != nil {
		return fmt.Errorf("TxService.SendCoinByName: %w", err)
	}

	return nil
}

func (s *Transaction) BuyItemByName(ctx context.Context, uid domain.UserID, name domain.MerchName) error {
	merch, err := s.merchRepo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("TxService.BuyItemByName: error while searching item: %w", err)
	}

	err = s.repo.BuyItem(ctx, uid, merch)
	if err != nil {
		return fmt.Errorf("TxService.BuyItemByName: error while making purchase: %w", err)
	}

	return nil
}
