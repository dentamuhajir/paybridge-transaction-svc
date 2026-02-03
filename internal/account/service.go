package account

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateAccountWithInitialBalances(ctx context.Context, userID uuid.UUID) error
	GetAccount(ctx context.Context, userID uuid.UUID) (Account,error)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) Service {
	return &service{repo, log}
}

func (s *service) CreateAccountWithInitialBalances(ctx context.Context, userID uuid.UUID) error {
	return s.repo.CreateAccountWithBalance(ctx, userID)
}

func (s *service) GetAccount(ctx context.Context, userID uuid.UUID) (Account,error) {
	return s.repo.GetAccount(ctx, userID)
}
