package props

import (
	"errors"
	"time"
)

type SaveCardReq struct {
	UserUUID      string    `json:"-"`
	CardNumber    string    `json:"card_number"`
	Date          time.Time `json:"date"`
	CVC           string    `json:"cvc"`
	PaymentSystem string    `json:"payment_system"`
	IsPreferred   bool      `json:"is_preferred"`
	Email         string    `json:"email"`
}

type SaveCardResp struct {
	Message string `json:"message"`
}

func (r SaveCardReq) Validate() error {
	if r.CardNumber == "" {
		return errors.New("card number is required")
	}
	if r.CVC == "" {
		return errors.New("cvc code is required")
	}
	if r.PaymentSystem == "" {
		return errors.New("payment system is required")
	}
	if r.CVC != "" {
		if len(r.CVC) > 3 {
			return errors.New("cvc code is invalid")
		}
	}
	if r.Date.Unix() < time.Now().UTC().Unix() {
		return errors.New("date is invalid")
	}
	return nil
}
