package loan

import (
	"context"
	"paybridge-transaction-service/internal/domain/loan/dto"
	"paybridge-transaction-service/internal/domain/loan/entity"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, req dto.LoanAppCreateReq) (*dto.LoanAppCreateResp, error)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, req dto.LoanAppCreateReq) (*dto.LoanAppCreateResp, error) {

	loan := entity.LoanApplication{
		UserID:       req.UserID,
		ProductID:    req.ProductID,
		Amount:       req.Amount,
		TenorMonth:   req.TenorMonth,
		InterestType: req.InterestType,
		AdminFee:     req.AdminFee,
	}

	result, err := s.repo.Create(ctx, loan)
	if err != nil {
		return nil, err
	}

	return &dto.LoanAppCreateResp{
		ID:     result.ID.String(),
		Status: result.Status,
	}, nil
}
