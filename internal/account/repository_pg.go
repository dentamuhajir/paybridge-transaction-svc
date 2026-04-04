package account

import (
	"context"
	"paybridge-transaction-service/internal/infra/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	span := trace.SpanFromContext(ctx)
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

			span.RecordError(err)
			span.SetStatus(codes.Error, "account not found")

			return Account{}, ErrAccountNotFound
		}
		r.log.Error(ctx, "failed to get account", err,
			zap.String("owner_id", ownerID.String()),
		)
		return Account{}, err
	}

	return account, nil
}

func (r *repository) CreateAccountTx(
	ctx context.Context,
	tx pgx.Tx,
	acc Account,
) (Account, error) {

	query := `
        INSERT INTO accounts (
            owner_type,
            owner_id,
            account_code,
            currency,
            reference_type,
            reference_id,
            status
        )
        VALUES ($1,$2,$3,$4,$5,$6,'ACTIVE')
        ON CONFLICT (owner_type, owner_id, account_code, currency, reference_id)
        DO UPDATE SET updated_at = NOW()
        RETURNING id, owner_type, owner_id, account_code,
                  currency, reference_type, reference_id,
                  status, created_at, updated_at
    `

	var created Account

	err := tx.QueryRow(ctx, query,
		acc.OwnerType,
		acc.OwnerID,
		acc.AccountCode,
		acc.Currency,
		acc.ReferenceType,
		acc.ReferenceID,
	).Scan(
		&created.ID,
		&created.OwnerType,
		&created.OwnerID,
		&created.AccountCode,
		&created.Currency,
		&created.ReferenceType,
		&created.ReferenceID,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return Account{}, err
	}

	return created, nil
}

func (r *repository) GetAccountBalance(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	span := trace.SpanFromContext(ctx)

	// Get the account first
	account, err := r.GetAccount(ctx, ownerID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "account not found")
		return 0, err
	}

	// Query balance from account_balances table
	var balance int64
	query := `
		SELECT balance
		FROM account_balances
		WHERE account_id = $1
	`
	err = r.db.QueryRow(ctx, query, account.ID).Scan(&balance)
	if err != nil {
		r.log.Error(ctx, "failed to get account balance", err,
			zap.String("account_id", account.ID.String()),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get account balance")
		return 0, err
	}
	return balance, nil
}
