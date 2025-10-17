package models

import (
	"github.com/google/uuid"
	"time"
)

type VerificationCode struct {
	UUID      uuid.UUID `json:"uuid"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}
