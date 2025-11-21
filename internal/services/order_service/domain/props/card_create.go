package props

import (
	"errors"
	card_types "k071123/internal/services/order_service/domain/models/card_types"
	"regexp"
)

type CreateCard struct {
	Last4         string              `json:"last4"`
	PaymentSystem card_types.CardType `json:"payment_system"`
	YookassaToken *string             `json:"yookassa_token"`
	Holder        string              `json:"holder"`
}

func (c *CreateCard) Validate() error {
	if len(c.Last4) != 4 {
		return errors.New("last4 must be 4 digits")
	}
	if match, _ := regexp.MatchString(`^\d{4}$`, c.Last4); !match {
		return errors.New("last4 must contain only digits")
	}

	validPaymentSystems := []card_types.CardType{
		card_types.MasterCard,
		card_types.VisaCard,
		card_types.Mir,
		card_types.UnionPay,
		card_types.JCB,
		card_types.AmericanExpress,
		card_types.DinersClub,
		card_types.DiscoverCard,
		card_types.InstaPayment,
		card_types.InstaPaymentTM,
		card_types.Laser,
		card_types.Dankort,
		card_types.Solomon,
		card_types.Switch,
	}
	isValidSystem := false
	for _, system := range validPaymentSystems {
		if c.PaymentSystem == system {
			isValidSystem = true
			break
		}
	}
	if !isValidSystem {
		return errors.New("invalid payment system, must be one of: MasterCard, Visa, Mir, UnionPay, JCB, AmericanExpress, DinersClub, DiscoverCard, InstaPayment, InstaPaymentTM, Laser, Dankort, Solomon, Switch")
	}

	if c.Holder == "" {
		return errors.New("card holder name is required")
	}

	return nil
}
