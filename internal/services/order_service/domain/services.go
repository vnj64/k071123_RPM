package domain

import (
	"k071123/internal/services/order_service/domain/services"
	"k071123/tools/logger"
)

type Services interface {
	Config() services.Config
	Billing() services.Billing
	Logger() *logger.Logger
}
