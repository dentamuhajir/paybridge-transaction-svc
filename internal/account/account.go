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
