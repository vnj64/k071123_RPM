package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) Add(model *models.Car) error {
	return r.db.Create(model).Error
}

func (r *CarRepository) GetByUUID(uuid string) (*models.Car, error) {
	var car models.Car
	if err := r.db.Where("uuid = ?", uuid).Preload("Settings").First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *CarRepository) GetByGosNumber(number string) (*models.Car, error) {
	var car models.Car
	if err := r.db.Where("gos_number = ?", number).Preload("Settings").First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}
