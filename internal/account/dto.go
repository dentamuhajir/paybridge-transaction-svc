package account

import "github.com/google/uuid"

type AccountResponse struct {
	OwnerID *uuid.UUID `json:"owner_id"`
	Status  Status     `json:"status"`
}
