package wallet

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, w *Wallet) (*Wallet, error)
}

type repository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewRepository(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repository{db: db, log: log}
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
