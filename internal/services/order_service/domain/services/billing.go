package services

import "time"

type Billing interface {
	GeneratePayToken(last4, cvc string, date time.Time) string
	AutoPay(token string, amount float64) error
}
