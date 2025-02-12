package postgres

import (
	"avito_shop/internal/domain"
	"avito_shop/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchRepository struct {
	pool *pgxpool.Pool
}

func NewMerchRepository(dbPool *pgxpool.Pool) repository.Merch {
	return &MerchRepository{
		pool: dbPool,
	}
}

func (r *MerchRepository) GetByName(ctx context.Context, name string) (domain.Merch, error) {
	item := domain.Merch{Name: name}

	query := `SELECT id, price
              FROM merch
              WHERE name = $1`

	err := r.pool.QueryRow(ctx, query, name).Scan(&item.ID, &item.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Merch{}, domain.ErrMerchNotFound
		}
		return domain.Merch{}, domain.ErrMerchNotFound
	}

	return item, nil
}
