package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository interface {
	CreateAccountWithBalance(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewRepository(db *pgxpool.Pool, log *zap.Logger) Repository {
	return &repository{db, log}
}

func (r *repository) CreateAccountWithBalance(ctx context.Context, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var accountID uuid.UUID

	err = tx.QueryRow(ctx, ` INSERT INTO accounts (owner_id, status) 
			VALUES ($1, 'ACTIVE') ON CONFLICT (owner_id) 
			DO NOTHING RETURNING id
		`, userID).Scan(&accountID)

	if err == pgx.ErrNoRows {
		err = tx.QueryRow(ctx, `
			SELECT id FROM accounts WHERE owner_id = $1
		`, userID).Scan(&accountID)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
	INSERT INTO account_balances (account_id, balance_type_id, amount)
	VALUES
		($1, 1, 0), -- CASH
		($1, 2, 0), -- LOAN_PRINCIPAL
		($1, 3, 0), -- LOAN_INTEREST
		($1, 4, 0), -- FEE
		($1, 5, 0), -- ESCROW
		($1, 6, 0)  -- RESERVE
	ON CONFLICT DO NOTHING
`, accountID)

	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}
