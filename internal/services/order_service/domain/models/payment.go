package models

import (
	"github.com/google/uuid"
	"k071123/internal/services/order_service/domain/models/payment_statuses"
	"k071123/pkg/timestamps"
)

type Payment struct {
	UUID        uuid.UUID `json:"uuid" gorm:"primaryKey"`
	SessionUUID uuid.UUID `json:"session_uuid"`
	//PaymentMethod payment_methods.PaymentMethod    `json:"payment_method"`
	Status        payment_statuses.PaymentStatuses `json:"status"`
	TransactionId string                           `json:"transaction_id"`
	Amount        float64                          `json:"amount"`
	PlatformFee   *float64                         `json:"platform_fee"`
	Description   string                           `json:"description"`
	Timestamps    timestamps.Timestamps            `gorm:"embedded" json:"timestamps" swaggerignore:"true"`
}
