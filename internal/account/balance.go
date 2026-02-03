package account

import (
	"time"

	"github.com/google/uuid"
)

type AccountBalance struct {
	AccountID     uuid.UUID `db:"account_id"`
	BalanceTypeID int32     `db:"balance_type_id"`
	Amount        int64     `db:"amount"`
	UpdatedAt     time.Time `db:"updated_at"`
}
