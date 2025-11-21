package repositories

import "k071123/internal/services/parking_service/domain/models"

type UnitRepository interface {
	Add(model *models.Unit) error
	GetByUUID(uuid string) (*models.Unit, error)
}
