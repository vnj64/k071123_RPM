package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
)

type ParkingRepository struct {
	db *gorm.DB
}

func NewParkingRepository(db *gorm.DB) *ParkingRepository {
	return &ParkingRepository{
		db: db,
	}
}

func (r *ParkingRepository) Add(model *models.Parking) error {
	return r.db.Create(model).Error
}

func (r *ParkingRepository) GetByUUID(uuid string) (*models.Parking, error) {
	var parking models.Parking
	if err := r.db.Where("uuid = ?", uuid).First(&parking).Error; err != nil {
		return nil, err
	}
	return &parking, nil
}

func (r *ParkingRepository) Update(model *models.Parking) error {
	return r.db.Save(model).Error
}
