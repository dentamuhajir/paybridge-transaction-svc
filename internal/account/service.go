package account

import (
	"context"
	"paybridge-transaction-service/internal/infra/logger"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	span := trace.SpanFromContext(ctx)

	if userID == uuid.Nil {
		err := ErrInvalidUserID
		s.log.Warn(ctx, "nil user id")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Account{}, ErrInvalidUserID
	}

	account, err := s.repo.GetAccount(ctx, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Account{}, err
	}

	if account.Status != StatusActive {
		err := ErrAccountInactive

		s.log.Info(ctx, "account inactive")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return Account{}, ErrAccountInactive
	}

	return account, err
}
