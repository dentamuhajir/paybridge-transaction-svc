package loan

import (
	"time"

	"github.com/google/uuid"
)

type LoanAppCreateRequest struct {
	UserID                  uuid.UUID `json:"user_id"`
	ProductID               uuid.UUID `json:"product_id"`
	Amount                  int64     `json:"amount"`
	TenorMonth              int       `json:"tenor_month"`
	InterestType            string    `json:"interest_type"` // FLAT / ANNUITY
	DisbursementScheduledAt time.Time `json:"disbursement_scheduled_at"`
	AdminFee                int64     `json:"admin_fee"`
}

type LoanAppCreateResponse struct {
	ID     string `db:"id"`
	Status string `db:"status"`
}

type LoanApprovalRequest struct {
	ID     uuid.UUID `json:"user_id"`
	Status string    `json:"status"`
}

type LoanApprovalResponse struct {
	ID     string `db:"id"`
	Status string `db:"status"`
}
