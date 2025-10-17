package models

import (
	"github.com/google/uuid"
	"k071123/internal/services/user_service/domain/models/user_status"
	"k071123/internal/shared/permissions"
	"k071123/pkg/timestamps"
	"time"
)

type User struct {
	Timestamps  timestamps.Timestamps `gorm:"embedded;" json:"timestamps" swaggerignore:"false"`
	UUID        uuid.UUID             `json:"uuid"`
	FirstName   *string               `json:"first_name"`
	SecondName  *string               `json:"second_name"`
	BirthDate   *time.Time            `json:"birth_date"`
	Status      user_status.Status    `json:"status"`
	Email       string                `json:"email"`
	Role        permissions.Role      `json:"role"`
	PhoneNumber *string               `json:"phone_number"`
}
