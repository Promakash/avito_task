package postgres

import (
	"avito_shop/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"avito_shop/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(dbPool *pgxpool.Pool) repository.Transaction {
	return &TransactionRepository{
		pool: dbPool,
	}
}

func (r *TransactionRepository) SendCoin(ctx context.Context, tx domain.Transaction) error {
	dbTx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer func() {
		if err != nil {
			_ = dbTx.Rollback(ctx)
		}
	}()

	err = r.updateUserBalance(ctx, dbTx, tx.From, -tx.Amount)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == PgCheckViolation {
				return fmt.Errorf("TxRepository.SendCoin: %w", domain.ErrLowBalance)
			}
		}
		return fmt.Errorf("TxRepository.SendCoin: %w", err)
	}

	err = r.updateUserBalance(ctx, dbTx, tx.To, tx.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("TxRepository.SendCoin: %w", domain.ErrUserNotFound)
		}
		return fmt.Errorf("TxRepository.SendCoin: %w", err)
	}

	err = r.insertTransaction(ctx, dbTx, tx)
	if err != nil {
		return fmt.Errorf("TxRepository.SendCoin: %w", err)
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("TxRepository.SendCoin: %w", err)
	}

	return nil
}

func (r *TransactionRepository) BuyItem(ctx context.Context, uid domain.UserID, item domain.Merch) error {
	dbTx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer func() {
		if err != nil {
			_ = dbTx.Rollback(ctx)
		}
	}()

	tx := domain.Transaction{
		From:   uid,
		To:     repository.ShopDBID,
		Amount: item.Price,
	}

	err = r.updateUserBalance(ctx, dbTx, tx.From, -tx.Amount)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == PgCheckViolation {
				return fmt.Errorf("TxRepository.BuyItem: %w", domain.ErrLowBalance)
			}
		}
		return fmt.Errorf("TxRepository.BuyItem: %w", err)
	}

	err = r.addItemToInventory(ctx, dbTx, tx.From, item)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("TxRepository.BuyItem: invalid inventory: %w", domain.ErrMerchNotFound)
		}
		return fmt.Errorf("TxRepository.BuyItem: %w", err)
	}

	err = r.insertTransaction(ctx, dbTx, tx)
	if err != nil {
		return fmt.Errorf("TxRepository.BuyItem: %w", err)
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("TxRepository.BuyItem: %w", err)
	}

	return nil
}

func (r *TransactionRepository) updateUserBalance(
	ctx context.Context,
	dbTx pgx.Tx,
	id domain.UserID,
	amount int,
) error {
	query := `UPDATE employees
              SET coins = coins + $2
              WHERE id = $1`

	_, err := dbTx.Exec(ctx, query, id, amount)
	if err != nil {
		return fmt.Errorf("TxRepository.updateUserBalance: %w", err)
	}
	return nil
}

func (r *TransactionRepository) insertTransaction(
	ctx context.Context,
	dbTx pgx.Tx,
	tx domain.Transaction,
) error {
	query := `INSERT INTO coin_transactions 
    		  (sender, recipient, amount)
              VALUES ($1, $2, $3)`

	_, err := dbTx.Exec(ctx, query, tx.From, tx.To, tx.Amount)
	if err != nil {
		return fmt.Errorf("TxRepository.insertTransaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) addItemToInventory(
	ctx context.Context,
	dbTx pgx.Tx,
	uid domain.UserID,
	item domain.Merch,
) error {
	query := `INSERT INTO inventory (employee_id, merch_id, quantity)
              VALUES ($1, $2, 1)
              ON CONFLICT (employee_id, merch_id) 
              DO UPDATE SET quantity = inventory.quantity + 1`

	_, err := dbTx.Exec(ctx, query, uid, item.ID)
	if err != nil {
		return fmt.Errorf("TxRepository.addItemToInventory: %w", err)
	}
	return nil
}
