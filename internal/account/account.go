package account

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type Account struct {
	ID            uuid.UUID
	OwnerType     string
	OwnerID       *uuid.UUID
	AccountCode   string
	Currency      string
	ReferenceType *string
	ReferenceID   *uuid.UUID
	Status        Status
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
