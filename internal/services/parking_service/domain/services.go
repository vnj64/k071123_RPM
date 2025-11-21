package domain

import "k071123/internal/services/parking_service/domain/services"

type Services interface {
	Config() services.Config
}
