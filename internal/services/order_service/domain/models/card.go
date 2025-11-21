package models

import (
	"github.com/google/uuid"
	"k071123/pkg/timestamps"
)

type Card struct {
	Timestamps    timestamps.Timestamps `gorm:"embedded" json:"-" swaggerignore:"true"`
	UUID          uuid.UUID             `json:"uuid" gorm:"primaryKey"`
	Last4Digits   string                `json:"last_4_digits" gorm:"column:last4"`
	PaymentSystem string                `json:"payment_system" gorm:"column:payment_system"`
	UserUUID      string                `json:"user_uuid" gorm:"column:user_uuid"`
	IsPreferred   bool                  `json:"is_preferred" gorm:"column:is_preferred"`
	Token         string                `json:"token"`
	IsActive      bool                  `json:"is_active"`
}
