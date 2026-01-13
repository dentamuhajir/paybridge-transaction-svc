package dto

import "github.com/google/uuid"

type LoanAppCreateReq struct {
	UserID       uuid.UUID `json:"user_id"`
	ProductID    uuid.UUID `json:"product_id"`
	Amount       int64     `json:"amount"`
	TenorMonth   int       `json:"tenor_month"`
	InterestType string    `json:"interest_type"` // FLAT / ANNUITY
	AdminFee     int64     `json:"admin_fee"`
}

type LoanAppCreateResp struct {
	ID     string `db:"id"`
	Status string `db:"status"`
}
