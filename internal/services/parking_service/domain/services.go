package domain

import (
	"k071123/internal/services/parking_service/domain/services"
	"k071123/tools/logger"
)

type Services interface {
	Config() services.Config
	Logger() *logger.Logger
}
