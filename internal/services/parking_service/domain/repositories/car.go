package repositories

import "k071123/internal/services/parking_service/domain/models"

type CarRepository interface {
	Add(model *models.Car) error
	GetByUUID(uuid string) (*models.Car, error)
	GetByGosNumber(number string) (*models.Car, error)
}
