package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
)

type TariffRepository struct {
	db *gorm.DB
}

func NewTariffRepository(db *gorm.DB) *TariffRepository {
	return &TariffRepository{
		db: db,
	}
}

func (r *TariffRepository) Add(model *models.Tariff) error {
	return r.db.Create(model).Error
}

func (r *TariffRepository) GetByUUID(uuid string) (*models.Tariff, error) {
	var tariff models.Tariff
	if err := r.db.Where("uuid = ?", uuid).First(&tariff).Error; err != nil {
		return nil, err
	}
	return &tariff, nil
}
