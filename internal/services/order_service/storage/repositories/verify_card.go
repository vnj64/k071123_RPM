package repositories

import (
	"errors"
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/repositories"
)

type verifyTokenRepository struct {
	db *gorm.DB
}

func NewVerifyTokenRepository(db *gorm.DB) repositories.VerifyTokenRepository {
	return &verifyTokenRepository{db: db}
}

func (r *verifyTokenRepository) Insert(model *models.VerifyTokens) error {
	return r.db.Create(model).Error
}

func (r *verifyTokenRepository) GetLastByUserUUID(userUuid string) (*models.VerifyTokens, error) {
	var model models.VerifyTokens
	if err := r.db.Where("user_uuid = ? AND used = false", userUuid).Order("created_at DESC").First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *verifyTokenRepository) Save(model *models.VerifyTokens) error {
	return r.db.Where("uuid = ?", model.UUID.String()).Save(model).Error
}
