package postgres

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"avito_shop/pkg/infra/cache"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool              *pgxpool.Pool
	cacheByName       cache.Cache
	cacheTTL          time.Duration
	cacheWriteTimeout time.Duration
}

func NewUserRepository(
	dbPool *pgxpool.Pool,
	cache cache.Cache,
	cacheTTL,
	cacheWriteTimeout time.Duration,
) repository.User {
	return &UserRepository{
		pool:              dbPool,
		cacheByName:       cache,
		cacheTTL:          cacheTTL,
		cacheWriteTimeout: cacheWriteTimeout,
	}
}

func (r *UserRepository) Put(ctx context.Context, user domain.User) (domain.UserID, error) {
	var id domain.UserID
	query := `INSERT INTO Employees (username, hashed_password)
              VALUES ($1, $2)
              RETURNING id`

	err := r.pool.QueryRow(ctx, query, user.Name, user.HashedPassword).Scan(&id)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == PgUniqueViolation {
				return 0, fmt.Errorf("UserRepository.Put: %w", domain.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("UserRepository.Put: %w", err)
	}

	return id, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name domain.UserName) (domain.User, error) {
	user := domain.User{}
	err := r.cacheByName.Get(ctx, name, &user)
	if err == nil {
		return user, nil
	}

	user.Name = name

	query := `SELECT id, hashed_password, coins FROM Employees
              WHERE username = $1`

	err = r.pool.QueryRow(ctx, query, name).Scan(&user.ID, &user.HashedPassword, &user.Info.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("UserRepository.GetByName: %w", domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("UserRepository.GetByName: %w", err)
	}

	go r.cacheUserByNameAsync(user)

	return user, nil
}

func (r *UserRepository) cacheUserByNameAsync(user domain.User) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cacheWriteTimeout)
	defer cancel()
	_ = r.cacheByName.Set(ctx, user.Name, user, r.cacheTTL)
}

func (r *UserRepository) GetInfoByID(ctx context.Context, id domain.UserID) (domain.UserInfo, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	coins, err := r.getCoinsTx(ctx, tx, id)
	if err != nil {
		return domain.UserInfo{}, fmt.Errorf("UserRepository.GetInfoByID: %w", err)
	}

	inv, err := r.getInventoryTx(ctx, tx, id)
	if err != nil {
		return domain.UserInfo{}, fmt.Errorf("UserRepository.GetInfoByID: %w", err)
	}

	txHistory, err := r.getTxHistoryTx(ctx, tx, id)
	if err != nil {
		return domain.UserInfo{}, fmt.Errorf("UserRepository.GetInfoByID: %w", err)
	}

	err = tx.Commit(ctx)

	return domain.UserInfo{
		Coins:        coins,
		Inventory:    inv,
		Transactions: txHistory,
	}, nil
}

func (r *UserRepository) getCoinsTx(ctx context.Context, tx pgx.Tx, id domain.UserID) (int, error) {
	var coins int

	query := `SELECT coins from Employees WHERE id = $1`
	err := tx.QueryRow(ctx, query, id).Scan(&coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("UserRepository.getCoinsTx: %w", domain.ErrUserNotFound)
		}
		return 0, fmt.Errorf("UserRepository.getCoinsTx: %w", err)
	}

	return coins, nil
}

func (r *UserRepository) getTxHistoryTx(
	ctx context.Context,
	tx pgx.Tx,
	id domain.UserID,
) ([]domain.UserTransaction, error) {
	query := `SELECT 
    		  ct.sender,
    		  COALESCE(e_from.username, 'deleted'),
    		  COALESCE(e_to.username, 'deleted'),
    		  ct.amount
			  FROM coin_transactions ct
              LEFT JOIN employees e_from ON ct.sender = e_from.id
              LEFT JOIN employees e_to   ON ct.recipient = e_to.id
              WHERE ct.sender = $1 OR ct.recipient = $1
              ORDER BY ct.created_at DESC`

	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UserRepository.getTxHistoryTx: %w", domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("UserRepository.getTxHistoryTx: %w", err)
	}
	defer func() { rows.Close() }()

	var txHistory []domain.UserTransaction
	for rows.Next() {
		var (
			idFrom   domain.UserID
			nameFrom = "deleted"
			nameTo   = "deleted"
			userTX   domain.UserTransaction
		)

		if err = rows.Scan(&idFrom, &nameFrom, &nameTo, &userTX.Amount); err != nil {
			return nil, fmt.Errorf("UserRepository.getTxHistoryTx: %w", err)
		}

		if idFrom == id {
			userTX.OtherUser = nameTo
			userTX.Direction = domain.Sent
		} else {
			userTX.OtherUser = nameFrom
			userTX.Direction = domain.Received
		}

		txHistory = append(txHistory, userTX)
	}

	return txHistory, nil
}

func (r *UserRepository) getInventoryTx(ctx context.Context, tx pgx.Tx, id domain.UserID) ([]domain.Inventory, error) {
	query := `SELECT 
    		  merch.name,
    		  inv.quantity
			  FROM inventory inv
              JOIN merch ON inv.merch_id = merch.id
              WHERE inv.employee_id = $1`

	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UserRepository.getInventoryTx: %w", domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("UserRepository.getInventoryTx: %w", err)
	}
	defer func() { rows.Close() }()

	var inv []domain.Inventory

	for rows.Next() {
		var curInv domain.Inventory

		if err = rows.Scan(&curInv.Name, &curInv.Quantity); err != nil {
			return nil, fmt.Errorf("UserRepository.getInventoryTx: %w", err)
		}

		inv = append(inv, curInv)
	}

	return inv, nil
}
