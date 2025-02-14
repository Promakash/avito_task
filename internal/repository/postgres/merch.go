package postgres

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type MerchRepository struct {
	pool        *pgxpool.Pool
	cacheByName sync.Map
}

func NewMerchRepository(dbPool *pgxpool.Pool) repository.Merch {
	return &MerchRepository{
		pool:        dbPool,
		cacheByName: sync.Map{},
	}
}

func (r *MerchRepository) GetByName(ctx context.Context, name string) (domain.Merch, error) {
	if val, ok := r.cacheByName.Load(name); ok {
		if merch, ok := val.(domain.Merch); ok {
			return merch, nil
		}
	}

	item := domain.Merch{Name: name}

	query := `SELECT id, price
              FROM merch
              WHERE name = $1`

	err := r.pool.QueryRow(ctx, query, name).Scan(&item.ID, &item.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Merch{}, fmt.Errorf("MerchRepository.GetByName: %w", domain.ErrMerchNotFound)
		}
		return domain.Merch{}, fmt.Errorf("MerchRepository.GetByName: %w", err)
	}

	r.cacheByName.Store(name, item)

	return item, nil
}
