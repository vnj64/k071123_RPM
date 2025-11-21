package domain

import "k071123/internal/services/order_service/domain/services"

type Services interface {
	Config() services.Config
	Billing() services.Billing
}
