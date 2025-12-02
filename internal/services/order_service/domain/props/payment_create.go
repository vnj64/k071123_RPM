package props

import (
	"errors"
	"github.com/google/uuid"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/models/payment_methods"
)

type CreatePaymentReq struct {
	SessionUUID   uuid.UUID                     `json:"session_uuid"`
	PaymentMethod payment_methods.PaymentMethod `json:"payment_method"`
	Amount        float64                       `json:"amount"`
	Description   string                        `json:"description"`
	UserUUID      *uuid.UUID                    `json:"user_uuid"`
	CardUUID      *uuid.UUID                    `json:"card_uuid"`
}

type CreatePaymentResp struct {
	Payment *models.Payment `json:"payment"`
}

func (p *CreatePaymentReq) Validate() error {
	if p.SessionUUID == uuid.Nil {
		return errors.New("session_uuid is required and cannot be empty")
	}

	if p.Amount < 0 {
		return errors.New("amount must be greater than zero")
	}

	if p.Description == "" {
		return errors.New("description is required")
	}

	return nil
}
