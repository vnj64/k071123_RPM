package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	RefreshTokenUUID uuid.UUID `json:"refresh_token_uuid"`
	UserUUID         uuid.UUID `json:"user_uuid"`
	ExpiresAt        time.Time `json:"expires_at"`
}
