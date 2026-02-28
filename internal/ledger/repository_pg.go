package ledger

import (
	"context"
	"paybridge-transaction-service/internal/infra/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

func NewRepository(db *pgxpool.Pool, log *logger.Logger) Repository {
	return &repository{db, log}
}

type Repository interface {
	InitializeBalanceTx(ctx context.Context, tx pgx.Tx, accountID uuid.UUID) error
}

func (r *repository) InitializeBalanceTx(
	ctx context.Context,
	tx pgx.Tx,
	accountID uuid.UUID,
) error {
	query := `
        INSERT INTO account_balances (account_id, balance)
        VALUES ($1, 0)
        ON CONFLICT (account_id) DO NOTHING
    `

	_, err := tx.Exec(ctx, query, accountID)
	return err
}
