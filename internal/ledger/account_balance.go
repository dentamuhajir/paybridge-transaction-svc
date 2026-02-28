package ledger

import (
	"time"

	"github.com/google/uuid"
)

type BalanceTypeID int32

type Balance struct {
	AccountID uuid.UUID `db:"account_id"`
	Balance   int64     `db:"balance"`
	UpdatedAt time.Time `db:"updated_at"`
}
