package domain

import "k071123/internal/services/notification_service/domain/services"

type Services interface {
	Config() services.Config
	Smtp() services.Smtp
}
