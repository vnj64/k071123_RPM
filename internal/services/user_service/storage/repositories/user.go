package repositories

import (
	"errors"
	"gorm.io/gorm"
	"k071123/internal/services/user_service/domain/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Add(model *models.User) error {
	return r.db.Create(model).Error
}

func (r *UserRepository) GetByUUID(uuid string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Save(user *models.User) error {
	return r.db.Where("uuid = ?", user.UUID).Save(user).Error
}
