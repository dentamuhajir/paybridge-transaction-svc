package account

import (
	"context"
	"paybridge-transaction-service/internal/infra/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type repository struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

func NewRepository(db *pgxpool.Pool, log *logger.Logger) Repository {
	return &repository{db, log}
}

func (r *repository) GetAccount(ctx context.Context, ownerID uuid.UUID) (Account, error) {

	var account Account
	query := `
		SELECT id, owner_id, status
		FROM accounts
		WHERE owner_id = $1
	`
	err := r.db.QueryRow(ctx, query, ownerID).Scan(
		&account.ID,
		&account.OwnerID,
		&account.Status,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.log.Warn(ctx, "account not found",
				zap.String("owner_id", ownerID.String()),
			)

			return Account{}, ErrAccountNotFound
		}
		r.log.Error(ctx, "failed to get account", err,
			zap.String("owner_id", ownerID.String()),
		)
		return Account{}, err
	}

	return account, nil
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
			VALUES ($1, $2) ON CONFLICT (owner_id) 
			DO NOTHING RETURNING id `, userID, StatusActive).Scan(&accountID)

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
		($1, $2, 0), 
		($1, $3, 0), 
		($1, $4, 0), 
		($1, $5, 0), 
		($1, $6, 0), 
		($1, $7, 0)
	ON CONFLICT DO NOTHING`,
		accountID,
		BalanceTypeCash,
		BalanceTypeLoanPrincipal,
		BalanceTypeLoanInterest,
		BalanceTypeFee,
		BalanceTypeEscrow,
		BalanceTypeReserve)

	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}
