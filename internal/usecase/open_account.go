package usecase

import (
	"context"

	"paybridge-transaction-service/internal/account"
	"paybridge-transaction-service/internal/ledger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OpenAccountUsecase struct {
	db          *pgxpool.Pool
	accountRepo account.Repository
	ledgerRepo  ledger.Repository
}

func NewOpenAccountUsecase(
	db *pgxpool.Pool,
	accountRepo account.Repository,
	ledgerRepo ledger.Repository,
) *OpenAccountUsecase {
	return &OpenAccountUsecase{
		db:          db,
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
	}
}

func (u *OpenAccountUsecase) ExecuteUserAccount(
	ctx context.Context,
	userID uuid.UUID,
) error {

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	acc := account.Account{
		OwnerType:   "USER",
		OwnerID:     &userID,
		AccountCode: "MAIN",
		Currency:    "IDR",
		Status:      "ACTIVE",
	}

	created, err := u.accountRepo.CreateAccountTx(ctx, tx, acc)
	if err != nil {
		return err
	}

	if err = u.ledgerRepo.InitializeBalanceTx(ctx, tx, created.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
