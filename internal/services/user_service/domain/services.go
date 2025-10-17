package domain

import "k071123/internal/services/user_service/domain/services"

type Services interface {
	Config() services.Config
}
