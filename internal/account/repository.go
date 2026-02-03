package account

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAccountWithBalance(ctx context.Context, userID uuid.UUID) error
	GetAccount(ctx context.Context, userID uuid.UUID) (Account, error)
}
