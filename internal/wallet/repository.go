package wallet

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, w *Wallet) (*Wallet, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, w *Wallet) (*Wallet, error) {
	query := `
		INSERT INTO wallets (user_id, currency)
		VALUES ($1, $2)
		RETURNING id, balance, status, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		w.UserID,
		w.Currency,
	).Scan(
		&w.ID,
		&w.Balance,
		&w.Status,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return w, nil

}
