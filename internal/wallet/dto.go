package wallet

type CreateWalletReq struct {
	UserID   string `json:"user_id" validate:"required,uuid4"`
	Currency string `json:"currency" validate:"required"`
}

type CreateWalletResp struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
	Status  string `json:"status"`
}
