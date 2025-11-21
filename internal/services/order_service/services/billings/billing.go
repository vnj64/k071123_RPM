package billings

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

type billingClient struct{}

func NewBillingClient() *billingClient {
	//
	return &billingClient{}
}

func hashHmacSha256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// TODO: по хорошему реализовать expDate, чтобы токен обновлялся раз в N времени
func (c *billingClient) GeneratePayToken(last4, cvc string, date time.Time) string {
	token := hashHmacSha256(last4+date.String(), cvc)
	return token
}

func (c *billingClient) AutoPay(token string, amount float64) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}
	if amount < 0 || amount == 0 {
		return errors.New("amount must be more than 0")
	}

	return nil
}
