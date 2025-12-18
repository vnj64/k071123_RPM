package props

import (
	"errors"
	"github.com/google/uuid"
)

type FinishSessionRequest struct {
	CarNumber     string        `json:"car_number"`
	UnitUUID      uuid.UUID     `json:"unit_uuid"`
	UserUUID      string        `json:"-"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

func (r FinishSessionRequest) Validate() error {
	if r.CarNumber == "" {
		return errors.New("car number is required")
	}
	// TODO: позже поправить, сделать проверку на typeof(PaymentMethod)
	if r.PaymentMethod != "bank_card" {
		return errors.New("payment method is required")
	}
	if r.UserUUID == "" {
		return errors.New("user uuid is required")
	}
	if r.UnitUUID == uuid.Nil {
		return errors.New("unit uuid is required")
	}
	return nil
}

type PaymentMethod string

const (
	BankCard PaymentMethod = "bank_card"
)

type FinishSessionResp struct {
	Status string `json:"status"`
}
