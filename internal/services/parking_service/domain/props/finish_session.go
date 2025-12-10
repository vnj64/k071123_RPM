package props

import "github.com/google/uuid"

type FinishSessionRequest struct {
	CarNumber     string        `json:"car_number"`
	UnitUUID      uuid.UUID     `json:"unit_uuid"`
	UserUUID      string        `json:"-"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type PaymentMethod string

const (
	BankCard PaymentMethod = "bank_card"
)

type FinishSessionResp struct {
	Status string `json:"status"`
}
