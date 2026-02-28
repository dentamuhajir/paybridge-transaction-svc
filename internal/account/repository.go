package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetAccount(ctx context.Context, userID uuid.UUID) (Account, error)
	CreateAccountTx(ctx context.Context, tx pgx.Tx, acc Account) (Account, error)
}
