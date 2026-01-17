package loan

import (
	"context"
	"paybridge-transaction-service/internal/loan/entity"
	"time"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, req LoanAppCreateRequest) (*LoanAppCreateResponse, error)
	Approval(ctx context.Context, req LoanApprovalRequest) (*LoanApprovalResponse, error)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, req LoanAppCreateRequest) (*LoanAppCreateResponse, error) {

	loan := entity.LoanApplication{
		UserID:                  req.UserID,
		ProductID:               req.ProductID,
		Amount:                  req.Amount,
		TenorMonth:              req.TenorMonth,
		InterestType:            req.InterestType,
		AdminFee:                req.AdminFee,
		DisbursementScheduledAt: req.DisbursementScheduledAt,
	}

	result, err := s.repo.Create(ctx, loan)
	if err != nil {
		log.Error(ctx, "error in service", err)
		return nil, err
	}

	return &LoanAppCreateResponse{
		ID:     result.ID.String(),
		Status: result.Status,
	}, nil
}

func (s *service) Approval(ctx context.Context, req LoanApprovalRequest) (*LoanApprovalResponse, error) {
	loan := entity.LoanApplication{
		ID:        req.ID,
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}

	result, err := s.repo.Approval(ctx, loan)

	if err != nil {
		return nil, err
	}

	return &LoanApprovalResponse{
		ID:     result.ID.String(),
		Status: req.Status,
	}, nil

}
