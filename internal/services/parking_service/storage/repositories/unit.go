package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

func (r *UnitRepository) Add(model *models.Unit) error {
	return r.db.Create(model).Error
}

func (r *UnitRepository) GetByUUID(uuid string) (*models.Unit, error) {
	var Unit models.Unit
	if err := r.db.Where("uuid = ?", uuid).First(&Unit).Error; err != nil {
		return nil, err
	}
	return &Unit, nil
}
