package models

import (
	"github.com/google/uuid"
	"k071123/pkg/timestamps"
)

type Car struct {
	Timestamps   timestamps.Timestamps `gorm:"embedded;" json:"-" swaggerignore:"true"`
	UUID         uuid.UUID             `json:"uuid"`
	UserUUID     uuid.UUID             `gorm:"index;" json:"user_uuid"`
	GosNumber    string                `json:"gos_number"`
	IsActive     bool                  `json:"is_active"`
	SettingsUUID *uuid.UUID            `json:"settings_uuid"`
	Settings     CarSettings           `json:"settings" gorm:"foreignKey:SettingsUUID;references:UUID"`
}

type CarSettings struct {
	UUID uuid.UUID `json:"uuid"`
	Vin  string    `json:"vin"`
}
