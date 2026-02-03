package account

import "github.com/google/uuid"

type AccountRequest struct {
	OwnerID uuid.UUID `json:"user_id"`
}
