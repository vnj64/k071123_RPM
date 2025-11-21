package props

import (
	"errors"
	"github.com/google/uuid"
	"k071123/internal/services/order_service/domain/models/payment_methods"
)

type CreatePayment struct {
	SessionUUID     uuid.UUID                     `json:"session_uuid"`
	PaymentMethod   payment_methods.PaymentMethod `json:"payment_method"`
	Amount          float64                       `json:"amount"`
	PlatformFee     *float64                      `json:"platform_fee"`
	Description     string                        `json:"description"`
	PaymentMethodId *string                       `json:"payment_method_id"`
	UserUUID        *uuid.UUID                    `json:"user_uuid"`
	CardUUID        *uuid.UUID                    `json:"card_uuid"`
}

func (p *CreatePayment) Validate() error {
	if p.SessionUUID == uuid.Nil {
		return errors.New("session_uuid is required and cannot be empty")
	}

	if p.Amount < 0 {
		return errors.New("amount must be greater than zero")
	}

	if p.PlatformFee != nil && *p.PlatformFee < 0 {
		return errors.New("platform_fee must be greater than or equal to 0")
	}

	if p.Description == "" {
		return errors.New("description is required")
	}

	return nil
}
