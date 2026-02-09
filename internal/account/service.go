package account

import (
	"context"
	"paybridge-transaction-service/internal/infra/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateAccountWithInitialBalances(ctx context.Context, userID uuid.UUID) error
	GetAccount(ctx context.Context, userID uuid.UUID) (Account, error)
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) Service {
	return &service{repo, log}
}

func (s *service) CreateAccountWithInitialBalances(ctx context.Context, userID uuid.UUID) error {
	return s.repo.CreateAccountWithBalance(ctx, userID)
}

func (s *service) GetAccount(ctx context.Context, userID uuid.UUID) (Account, error) {
	if userID == uuid.Nil {
		s.log.Warn(ctx, "nil user id",
			zap.String("operation", "GetAccount"),
		)
		return Account{}, ErrInvalidUserID
	}

	account, err := s.repo.GetAccount(ctx, userID)

	if err != nil {
		return Account{}, err
	}

	if account.Status != StatusActive {
		s.log.Info(ctx, "account inactive",
			zap.String("owner_id", userID.String()),
			zap.String("status", string(account.Status)),
		)
		return Account{}, ErrAccountInactive
	}

	return account, err
}
