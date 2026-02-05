package account

import (
	"time"

	"github.com/google/uuid"
)

type BalanceTypeID int32

type Balance struct {
	AccountID uuid.UUID     `db:"account_id"`
	TypeID    BalanceTypeID `db:"balance_type_id"`
	Amount    int64         `db:"amount"`
	UpdatedAt time.Time     `db:"updated_at"`
}
