package ledger

import (
	"paybridge-transaction-service/internal/infra/logger"
)

type Service interface {
	//CreateAccountWithInitialBalances(ctx context.Context, userID uuid.UUID) error
	//GetAccount(ctx context.Context, userID uuid.UUID) (Account, error)
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) Service {
	return &service{repo, log}
}
