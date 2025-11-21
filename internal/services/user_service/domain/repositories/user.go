package repositories

import "k071123/internal/services/user_service/domain/models"

type User interface {
	Add(model *models.User) error
	GetByUUID(uuid string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}
