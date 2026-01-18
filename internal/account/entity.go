package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `db:"id"`
	OwnerID   uuid.UUID `db:"owner_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AccountBalance struct {
	AccountID     uuid.UUID `db:"account_id"`
	BalanceTypeID int32     `db:"balance_type_id"`
	Amount        int64     `db:"amount"`
	UpdatedAt     time.Time `db:"updated_at"`
}
