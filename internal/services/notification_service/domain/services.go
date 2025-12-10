package domain

import (
	"k071123/internal/services/notification_service/domain/services"
	"k071123/tools/logger"
)

type Services interface {
	Config() services.Config
	Smtp() services.Smtp
	Logger() *logger.Logger
	Amqp() services.AMQP
}
