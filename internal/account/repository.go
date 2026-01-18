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
	// Start a database transaction (BEGIN)
	// From this point on, all queries using `tx` are part of ONE atomic unit
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})

	committed := false
	defer func() {
		// If the transaction cannot be started (DB down, context canceled, etc.)
		// we must stop immediately — no work can be done safely
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Ensure cleanup happens no matter how the function exits
	defer func() {
		// If `err` is NOT nil, it means something failed after BEGIN
		// In that case, we MUST rollback to avoid partial writes
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var accountID uuid.UUID

	err = tx.QueryRow(ctx, ` INSERT INTO accounts (owner_id, status) 
			VALUES ($1, 'ACTIVE') ON CONFLICT (owner_id) 
			DO NOTHING RETURNING id
		`, userID).Scan(&accountID)

	// If no row was returned, it means the INSERT was skipped
	// because an account for this owner_id already exists.
	if err == pgx.ErrNoRows {
		// In this case, fetch the existing account ID
		// so we can continue initialization safely.
		err = tx.QueryRow(ctx, `
			SELECT id FROM accounts WHERE owner_id = $1
		`, userID).Scan(&accountID)

		// If we cannot fetch the existing account, this is a real error
		// and we must abort the transaction.
		if err != nil {
			return err
		}
		// If any other error occurred during INSERT or SELECT,
		// this is unexpected and must fail the transaction.
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

	// If balance initialization fails, we must abort
	// so the deferred rollback restores the DB state.
	if err != nil {
		return err
	}

	// ------------------------------------------------------------
	// 3️ Commit the transaction
	// ------------------------------------------------------------
	//
	// At this point:
	// - account exists
	// - all balances exist
	// - data is consistent
	//
	// Commit makes all changes permanent.

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}
