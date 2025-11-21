package repositories

import (
	"k071123/internal/services/parking_service/domain/models"
)

type ParkingFilter interface {
}

type ParkingUpdates interface {
}

type ParkingRepository interface {
	Add(model *models.Parking) error
	GetByUUID(uuid string) (*models.Parking, error)
}
