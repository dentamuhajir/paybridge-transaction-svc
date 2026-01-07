package wallet

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateWallet(ctx context.Context, req CreateWalletReq) (*CreateWalletResp, error)
	InquiryWallet(ctx context.Context, userID string) (*InquiryWalletResp, error)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) CreateWallet(ctx context.Context, req CreateWalletReq) (*CreateWalletResp, error) {

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		UserID:   userUUID,
		Balance:  0,
		Currency: req.Currency,
		Status:   "ACTIVE",
	}

	result, err := s.repo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return &CreateWalletResp{
		ID:      result.ID.String(),
		UserID:  result.UserID.String(),
		Balance: result.Balance,
		Status:  result.Status,
	}, nil

}

func (s *service) InquiryWallet(ctx context.Context, userID string) (*InquiryWalletResp, error) {

	result, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &InquiryWalletResp{
		UserID:  result.UserID.String(),
		Balance: result.Balance,
		Status:  result.Status,
	}, nil

}
