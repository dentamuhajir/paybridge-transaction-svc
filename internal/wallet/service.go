package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateWallet(ctx context.Context, req CreateWalletReq) (*CreateWalletResp, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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
