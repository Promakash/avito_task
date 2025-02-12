package postgres

import (
	"avito_shop/internal/domain"
	"context"
	"errors"
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

func (r *TransactionRepository) Put(ctx context.Context, tx domain.Transaction) error {
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
				return domain.ErrLowBalance
			}
		}
		return err
	}

	err = r.updateUserBalance(ctx, dbTx, tx.To, tx.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	err = r.insertTransaction(ctx, dbTx, tx)
	if err != nil {
		return err
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) PurchaseItem(ctx context.Context, tx domain.Transaction) error {
	dbTx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer func() {
		if err != nil {
			_ = dbTx.Rollback(ctx)
		}
	}()

	tx.To = repository.ShopDBID

	err = r.updateUserBalance(ctx, dbTx, tx.From, -tx.Amount)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == PgCheckViolation {
				return domain.ErrLowBalance
			}
		}
		return err
	}

	err = r.insertTransaction(ctx, dbTx, tx)
	if err != nil {
		return err
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
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
	return err
}

func (r *TransactionRepository) insertTransaction(
	ctx context.Context,
	dbTx pgx.Tx,
	tx domain.Transaction,
) error {
	query := `INSERT INTO coin_transactions 
    		  (from, to, amount)
              VALUES ($1, $2, $3)`

	_, err := dbTx.Exec(ctx, query, tx.From, tx.To, tx.Amount)
	return err
}
