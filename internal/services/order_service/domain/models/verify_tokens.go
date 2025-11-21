package models

import (
	"github.com/google/uuid"
	"time"
)

type VerifyTokens struct {
	UUID      uuid.UUID `json:"uuid"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Token     string    `json:"token"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}
