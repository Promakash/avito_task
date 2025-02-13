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
		return err
	}

	if toUser.ID == tx.From {
		return domain.ErrBadRequest
	}
	tx.To = toUser.ID

	return s.repo.SendCoin(ctx, tx)
}

func (s *Transaction) BuyItemByName(ctx context.Context, uid domain.UserID, name domain.MerchName) error {
	merch, err := s.merchRepo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("BuyItemByName: invalid item name: %w", domain.ErrBadRequest)
	}

	err = s.repo.BuyItem(ctx, uid, merch)
	if err != nil {
		return err
	}

	return nil
}
