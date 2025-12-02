package domain

import (
	"github.com/sirupsen/logrus"
	"k071123/internal/services/parking_service/domain/services"
)

type Services interface {
	Config() services.Config
	Logger() *logrus.Logger
}
