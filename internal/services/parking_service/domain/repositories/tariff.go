package repositories

import "k071123/internal/services/parking_service/domain/models"

type TariffRepository interface {
	Add(model *models.Tariff) error
	GetByUUID(uuid string) (*models.Tariff, error)
}
