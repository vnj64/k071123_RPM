package repositories

import (
	"gorm.io/gorm"
	"k071123/internal/services/user_service/domain/models"
)

type VerificationCodeRepository struct {
	db *gorm.DB
}

func NewVerificationCodeRepository(db *gorm.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{
		db: db,
	}
}

func (r *VerificationCodeRepository) GetByCode(code string) (*models.VerificationCode, error) {
	var model models.VerificationCode
	if err := r.db.Where("code = ? AND used = false", code).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *VerificationCodeRepository) GetByUUID(uuid string) (*models.VerificationCode, error) {
	var model models.VerificationCode
	if err := r.db.Where("uuid = ?", uuid).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *VerificationCodeRepository) GetByEmail(email string) (*models.VerificationCode, error) {
	var model models.VerificationCode
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *VerificationCodeRepository) Add(model *models.VerificationCode) error {
	return r.db.Create(model).Error
}

func (r *VerificationCodeRepository) GetLastByEmail(email string) (*models.VerificationCode, error) {
	var model models.VerificationCode
	if err := r.db.Where("email = ?", email).Order("created_at DESC").First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}
